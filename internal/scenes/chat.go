package scenes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/SaschaRunge/llm-local/internal/communication"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/SaschaRunge/llm-local/internal/llm"

	"github.com/google/uuid"
)

const (
	generateAnswerTimeoutInSeconds = 120
)

var SystemUserID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
var SystemUserName = "System"

type Chat struct {
	card           json.RawMessage
	userCharacter  character
	cachedMessages []cachedMessage
	characters     []character

	ID   uuid.UUID
	Name string

	runtime   core.Runtime
	llmClient *llm.Client
}

type cachedMessage struct {
	id       uuid.UUID
	authorID uuid.UUID
	communication.Message
}

type character = database.GetCharactersInChatRow

func NewSceneChat(runtime core.Runtime, chat database.Chat) (*Chat, error) {
	sceneChat := Chat{
		runtime: runtime,
		ID:      chat.ID,
		Name:    chat.Name,
	}

	if err := sceneChat.loadData(); err != nil {
		return nil, err
	}

	//TODO: do not hardcode model type/*
	llmClient, err := llm.NewClient("qwen3.6-27b")
	sceneChat.llmClient = llmClient

	if err != nil {
		return &Chat{}, err
	}

	return &sceneChat, nil
}

func (c *Chat) OnEnter() string {
	return ""
}

func (c *Chat) GetName() string {
	return fmt.Sprintf("%s", c.Name)
}

func (c *Chat) AtCharacter(name, userInput string) (core.Result, error) {
	var currentCharacter character
	for _, char := range c.characters {
		if char.Name == name {
			currentCharacter = char
			break
		}
	}

	if currentCharacter.Name == "" {
		return core.Result{}, fmt.Errorf("character %q not found in chat %q", name, c.GetName())
	}

	authorsNote := fmt.Sprintf("\n\n[SYSTEM: Respond strictly as %q for this turn. Maintain their tone and knowledge.]", currentCharacter.Name)

	history := append([]cachedMessage{}, c.cachedMessages...)
	userMessage := cachedMessage{
		authorID: c.userCharacter.ID,
		Message: communication.Message{
			Name:      c.userCharacter.Name,
			Role:      communication.RoleUser,
			Reasoning: "",
			Content:   userInput,
		},
	}
	alteredMessage := userMessage
	alteredMessage.Content = alteredMessage.Content + authorsNote

	history = append(history, alteredMessage)

	ctx, cancel := context.WithTimeout(c.runtime.Context(), generateAnswerTimeoutInSeconds*time.Second)
	defer cancel()

	response, err := c.llmClient.GenerateAnswer(ctx, asComMessages(history))
	if err != nil {
		return core.Result{}, err
	}

	assistantMessage := cachedMessage{
		authorID: currentCharacter.ID,
		Message: communication.Message{
			Name:      c.characters[0].Name,
			Role:      communication.RoleAssistant,
			Reasoning: "",
			Content:   response,
		},
	}

	err = c.archiveCurrentTurn(userMessage, assistantMessage)
	if err != nil {
		return core.Result{}, err
	}

	return core.Result{
		Response:  response + "\n",
		NextScene: c,
	}, nil
}

// TODO: need to add id/authorID to cached Messages after commiting to DB
// instead append to history at the end
func (c *Chat) HandleRawInput(userInput string) (core.Result, error) {
	if len(c.characters) == 0 {
		return core.Result{}, fmt.Errorf("no character in current chat %q", c.GetName())
	}

	history := append([]cachedMessage{}, c.cachedMessages...)
	userMessage := cachedMessage{
		authorID: c.userCharacter.ID,
		Message: communication.Message{
			Name:      c.userCharacter.Name,
			Role:      communication.RoleUser,
			Reasoning: "",
			Content:   userInput,
		},
	}

	history = append(history, userMessage)

	ctx, cancel := context.WithTimeout(c.runtime.Context(), generateAnswerTimeoutInSeconds*time.Second)
	defer cancel()
	response, err := c.llmClient.GenerateAnswer(ctx, asComMessages(history))
	if err != nil {
		return core.Result{}, err
	}

	assistantMessage := cachedMessage{
		//TODO: need to specify which character is talking. currently causes index out of range when
		//no characters in chat
		authorID: c.characters[0].ID,
		Message: communication.Message{
			Name:      c.characters[0].Name,
			Role:      communication.RoleAssistant,
			Reasoning: "",
			Content:   response,
		},
	}

	err = c.archiveCurrentTurn(userMessage, assistantMessage)
	if err != nil {
		return core.Result{}, err
	}

	return core.Result{
		Response:  response + "\n",
		NextScene: c,
	}, nil
}

func (c *Chat) archiveCurrentTurn(userMessage, assistantMessage cachedMessage) error {
	//TODO: rollback entire commit if send fails on either Add
	userMessageInDB, err := c.runtime.Store().DBQueries.AddMessage(c.runtime.Context(), database.AddMessageParams{
		Reasoning: sql.NullString{},
		Content:   userMessage.Content,
		ChatID:    c.ID,
		AuthorID:  c.userCharacter.ID,
		Role:      string(userMessage.Role),
	})
	if err != nil {
		return err
	}

	//TODO: currently answers just as the first character in the character list
	assistantMessageInDB, err := c.runtime.Store().DBQueries.AddMessage(c.runtime.Context(), database.AddMessageParams{
		Reasoning: sql.NullString{},
		Content:   assistantMessage.Content,
		ChatID:    c.ID,
		AuthorID:  assistantMessage.authorID,
		Role:      string(assistantMessage.Role),
	})
	if err != nil {
		return err
	}

	userMessage.id = userMessageInDB.ID
	assistantMessage.id = assistantMessageInDB.ID

	c.cachedMessages = append(c.cachedMessages, userMessage, assistantMessage)

	return nil
}

// do not allow system message to be regenerated
func (c *Chat) Regenerate(userInput string) (core.Result, error) {
	history := append([]cachedMessage{}, c.cachedMessages...)

	if len(history) == 0 {
		return core.Result{}, fmt.Errorf("no chat history, nothing to regenerate")
	}

	if history[len(history)-1].Role != communication.RoleAssistant {
		return core.Result{}, fmt.Errorf("can not regenerate user message")
	}

	lastMessage := history[len(history)-1]
	regenerationPrompt := fmt.Sprintf(
		"[SYSTEM: The user requested a regeneration of your previous answer with the following comment: %q. Please rephrase accordingly. Your previous answer was: %q.]",
		userInput, lastMessage.Content)
	message := cachedMessage{
		Message: communication.Message{
			Name:      SystemUserName,
			Role:      communication.RoleUser,
			Reasoning: "",
			Content:   regenerationPrompt,
		},
	}

	history[len(history)-1] = message

	ctxResponse, cancel := context.WithTimeout(c.runtime.Context(), generateAnswerTimeoutInSeconds*time.Second)
	defer cancel()
	response, err := c.llmClient.GenerateAnswer(ctxResponse, asComMessages(history))
	if err != nil {
		return core.Result{}, err
	}

	answer := cachedMessage{
		id:       lastMessage.id,
		authorID: lastMessage.authorID,
		Message: communication.Message{
			Name:      lastMessage.Name,
			Role:      communication.RoleAssistant,
			Reasoning: "",
			Content:   response,
		},
	}

	dbQueries := c.runtime.Store().DBQueries
	ctx := c.runtime.Context()

	err = dbQueries.ReplaceMessage(ctx, database.ReplaceMessageParams{
		ID:        answer.id,
		Reasoning: sql.NullString{},
		Content:   answer.Content,
	})
	if err != nil {
		return core.Result{}, err
	}

	c.cachedMessages[len(history)-1] = answer

	return core.Result{
		Response:  response,
		NextScene: c,
	}, nil
}

func (c *Chat) GetAvailableCharacters() []character {
	availableCharacters := []character{}
	for _, char := range c.characters {
		availableCharacters = append(availableCharacters, char)
	}
	return availableCharacters
}

func (c *Chat) loadData() error {
	db := c.runtime.Store().DBQueries
	ctx := c.runtime.Context()

	messages, err := db.GetChatHistory(ctx, c.ID)
	if err != nil {
		return err
	}

	c.characters, err = db.GetCharactersInChat(ctx, c.ID)
	if err != nil {
		return err
	}

	userCharacter, err := db.GetUserCharacterInChatByID(ctx, c.ID)
	if err != nil {
		return err
	}
	c.userCharacter = mapFromDBUserCharacter(userCharacter)

	card, err := db.GetCardFromChatByID(ctx, c.ID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(card, &c.card)
	if err != nil {
		return err
	}

	addToSystemPrompt := systemPromptBuilder()
	systemPrompt, err := addToSystemPrompt(c.card)
	if err != nil {
		return err
	}

	c.cachedMessages = append(c.cachedMessages, cachedMessage{
		Message: communication.Message{
			Name:      SystemUserName,
			Role:      communication.RoleSystem,
			Reasoning: "",
			Content:   systemPrompt,
		},
	})

	for _, message := range messages {
		role := communication.Role(message.Role)
		if !role.IsValid() {
			return fmt.Errorf("message role %q is not a valid role", role)
		}
		c.cachedMessages = append(c.cachedMessages, cachedMessage{
			id:       message.ID,
			authorID: message.AuthorID,
			Message: communication.Message{
				Name:      message.AuthorName,
				Role:      role,
				Reasoning: message.Reasoning.String,
				Content:   message.Content,
			},
		})
	}

	return nil
}

func asComMessages(messages []cachedMessage) []communication.Message {
	comMessages := []communication.Message{}
	for _, message := range messages {
		comMessages = append(comMessages, message.Message)
	}
	return comMessages
}

func systemPromptBuilder() func(json.RawMessage) (string, error) {
	var systemPrompt strings.Builder

	fmt.Fprintf(&systemPrompt, "%s\n", llm.SystemPrompt)

	return func(card json.RawMessage) (string, error) {
		cardStruct := core.Card{}
		err := json.Unmarshal(card, &cardStruct)

		if err != nil {
			return "", err
		}
		fmt.Fprintf(&systemPrompt, "%s\n", cardStruct.Description)
		return systemPrompt.String(), nil
	}
}

func mapFromDBUserCharacter(userCharacter database.GetUserCharacterInChatByIDRow) character {
	return character{
		ID:   userCharacter.ID,
		Name: userCharacter.Name,
		Card: userCharacter.Card,
	}
}
