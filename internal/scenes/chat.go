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
	generateAnswerTimeoutInSeconds = 480
)

var SystemUserID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
var SystemUserName = "System"

type Chat struct {
	card           json.RawMessage
	userCharacter  character
	cachedMessages []cachedMessage
	characters     []character
	systemPrompt   string

	lastResponseIdx int

	ID   uuid.UUID
	Name string

	runtime   core.Runtime
	llmClient *llm.Client
}

type cachedMessage struct {
	id       uuid.UUID
	authorID uuid.UUID
	variants []cachedMessage

	communication.Message
}

type character = database.GetCharactersInChatRow

func NewSceneChat(runtime core.Runtime, chat database.Chat) (*Chat, error) {
	sceneChat := Chat{
		runtime: runtime,
		ID:      chat.ID,
		Name:    chat.Name,
	}

	//TODO: do not hardcode model type/*
	llmClient, err := llm.NewClient("qwen3.6-27b")
	if err != nil {
		return &Chat{}, err
	}

	sceneChat.llmClient = llmClient

	if err := sceneChat.loadData(); err != nil {
		return nil, err
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

	history := append([]cachedMessage{}, c.cachedMessages[:len(c.cachedMessages)-1]...)
	lastResponse, err := c.getLastResponse()
	if err != nil {
		return core.Result{}, err
	}
	history = append(history, lastResponse)
	userMessage := cachedMessage{
		authorID: c.userCharacter.ID,
		Message: communication.Message{
			AuthorName: c.userCharacter.Name,
			Role:       communication.RoleUser,
			Reasoning:  "",
			Content:    userInput,
		},
	}
	alteredMessage := userMessage
	alteredMessage.Content = alteredMessage.Content + authorsNote

	history = append(history, alteredMessage)

	ctx, cancel := context.WithTimeout(c.runtime.Context(), generateAnswerTimeoutInSeconds*time.Second)
	defer cancel()

	reasoning, content, err := c.llmClient.GenerateAnswer(ctx, asComMessages(history))
	if err != nil {
		return core.Result{}, err
	}

	assistantMessage := cachedMessage{
		authorID: currentCharacter.ID,
		Message: communication.Message{
			AuthorName: currentCharacter.Name,
			Role:       communication.RoleAssistant,
			Reasoning:  reasoning,
			Content:    content,
		},
	}

	if userMessage.Content != "" {
		if err = c.archiveMessage(userMessage); err != nil {
			return core.Result{}, err
		}
	}
	if err = c.archiveMessage(assistantMessage); err != nil {
		return core.Result{}, err
	}

	c.lastResponseIdx = 0

	return core.Result{
		Author:    currentCharacter.Name,
		Response:  content,
		NextScene: c,
	}, nil
}

func (c *Chat) Cycle(idx int) (core.Result, error) {
	if len(c.cachedMessages) < 1 {
		return core.Result{}, fmt.Errorf("message history is empty, nothing to cycle through")
	}

	if c.cachedMessages[len(c.cachedMessages)-1].Role != communication.RoleAssistant {
		return core.Result{}, fmt.Errorf("unexpected: the last message in history is a user/system massage")
	}

	maxVariants := len(c.cachedMessages[len(c.cachedMessages)-1].variants)

	if maxVariants == 0 {
		return core.Result{}, fmt.Errorf("no variants of previous message available, nothing to cycle through")
	}

	if idx < 0 {
		c.lastResponseIdx = (c.lastResponseIdx + 1) % (maxVariants + 1)
	} else {
		if idx == 0 || idx > maxVariants+1 {
			return core.Result{}, fmt.Errorf("no variant with index %d. please choose an index between 1 and %d", idx, maxVariants+1)
		}
		c.lastResponseIdx = idx - 1
	}

	lastResponse, err := c.getLastResponse()
	if err != nil {
		return core.Result{}, nil
	}

	return core.Result{
		Author:    lastResponse.AuthorName,
		Response:  fmt.Sprintf("%s (%d/%d)", lastResponse.Content, c.lastResponseIdx+1, maxVariants+1),
		NextScene: c,
	}, nil
}

// TODO: need to add id/authorID to cached Messages after commiting to DB
// instead append to history at the end
func (c *Chat) HandleRawInput(userInput string) (core.Result, error) {
	if len(c.characters) == 0 {
		return core.Result{}, fmt.Errorf("no character in current chat %q", c.GetName())
	}

	history := append([]cachedMessage{}, c.cachedMessages[:len(c.cachedMessages)-1]...)
	lastResponse, err := c.getLastResponse()
	if err != nil {
		return core.Result{}, err
	}
	history = append(history, lastResponse)

	//overwrite cachedMessage with activeResponse. Still have to make sure the variants contain all messages
	c.cachedMessages[len(c.cachedMessages)-1].Content = lastResponse.Content
	c.cachedMessages[len(c.cachedMessages)-1].Reasoning = lastResponse.Reasoning

	userMessage := cachedMessage{
		authorID: c.userCharacter.ID,
		Message: communication.Message{
			AuthorName: c.userCharacter.Name,
			Role:       communication.RoleUser,
			Reasoning:  "",
			Content:    userInput,
		},
	}

	history = append(history, userMessage)

	/*for i, msg := range history {
		fmt.Printf("\nMessage in history: %d. %s\n", i, msg.Content)
	}*/

	ctx, cancel := context.WithTimeout(c.runtime.Context(), generateAnswerTimeoutInSeconds*time.Second)
	defer cancel()
	reasoning, content, err := c.llmClient.GenerateAnswer(ctx, asComMessages(history))
	if err != nil {
		return core.Result{}, err
	}

	assistantMessage := cachedMessage{
		//TODO: need to specify which character is talking. currently causes index out of range when
		//no characters in chat
		authorID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Message: communication.Message{
			AuthorName: SystemUserName,
			Role:       communication.RoleAssistant,
			Reasoning:  reasoning,
			Content:    content,
		},
	}

	err = c.archiveCurrentTurn(userMessage, assistantMessage)
	if err != nil {
		return core.Result{}, err
	}

	c.lastResponseIdx = 0

	return core.Result{
		Author:    SystemUserName,
		Response:  content,
		NextScene: c,
	}, nil
}

// do not allow system message to be regenerated
func (c *Chat) Regenerate(userInput string) (core.Result, error) {
	history := append([]cachedMessage{}, c.cachedMessages...)

	if len(history) <= 1 {
		return core.Result{}, fmt.Errorf("no chat history, nothing to regenerate")
	}

	lastResponse, err := c.getLastResponse()
	if err != nil {
		return core.Result{}, err
	}
	parentMessage := history[len(history)-1]

	if lastResponse.Role != communication.RoleAssistant {
		return core.Result{}, fmt.Errorf("can not regenerate a user or system message")
	}
	if lastResponse.Content == "" {
		return core.Result{}, fmt.Errorf("unexpected: last active response is empty")
	}

	if userInput != "" {
		regenerationPrompt := fmt.Sprintf(
			"[SYSTEM: The user requested a regeneration of your previous answer with the following comment: \"\"\"\n%s\"\"\"\n. Please craft a new response based on your previous response and the users critique. Your previous answer was: \"\"\"\n%s\"\"\"\n.]",
			userInput, lastResponse.Content)

		//fmt.Printf("TEST OUTPUT: %s\n\n\n", regenerationPrompt)
		message := cachedMessage{
			Message: communication.Message{
				AuthorName: SystemUserName,
				Role:       communication.RoleUser,
				Reasoning:  "",
				Content:    regenerationPrompt,
			},
		}

		history[len(history)-1] = message
	} else {
		history = append([]cachedMessage{}, history[:len(history)-1]...)

		if lastResponse.Role == communication.RoleAssistant && lastResponse.AuthorName != SystemUserName {
			authorsNote := fmt.Sprintf("\n\n[SYSTEM: Respond strictly as %q for this turn. Maintain their tone and knowledge.]", lastResponse.AuthorName)

			history = append(history, cachedMessage{
				authorID: c.userCharacter.ID,
				Message: communication.Message{
					AuthorName: c.userCharacter.Name,
					Role:       communication.RoleUser,
					Reasoning:  "",
					Content:    authorsNote,
				},
			})
		}
	}

	ctxResponse, cancel := context.WithTimeout(c.runtime.Context(), generateAnswerTimeoutInSeconds*time.Second)
	defer cancel()
	reasoning, content, err := c.llmClient.GenerateAnswer(ctxResponse, asComMessages(history))
	if err != nil {
		return core.Result{}, err
	}

	/*for i, msg := range history {
		fmt.Printf("\nMessage in history: %d. %s\n", i, msg.Content)
	}*/

	dbQueries := c.runtime.Store().DBQueries
	ctx := c.runtime.Context()

	variant, err := dbQueries.AddVariant(ctx, database.AddVariantParams{
		Reasoning:   sql.NullString{String: reasoning, Valid: true},
		Content:     content,
		ChatID:      c.ID,
		AuthorID:    lastResponse.authorID,
		Role:        string(lastResponse.Role),
		ParentMsgID: uuid.NullUUID{UUID: parentMessage.id, Valid: true},
	})
	if err != nil {
		return core.Result{}, err
	}

	newVariant := cachedMessage{
		id:       variant.ID,
		authorID: variant.AuthorID,
		Message: communication.Message{
			AuthorName: parentMessage.AuthorName,
			Role:       communication.RoleAssistant,
			Reasoning:  variant.Reasoning.String,
			Content:    variant.Content,
		},
	}

	c.cachedMessages[len(c.cachedMessages)-1].variants = append(c.cachedMessages[len(c.cachedMessages)-1].variants, newVariant)
	maxVariants := len(c.cachedMessages[len(c.cachedMessages)-1].variants)
	c.lastResponseIdx = maxVariants

	// fmt.Printf("\n\n %s", c.cachedMessages[c.lastResponseIdx].Content)

	return core.Result{
		Author:    lastResponse.AuthorName,
		Response:  fmt.Sprintf("%s (%d/%d)", content, maxVariants+1, maxVariants+1),
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

func (c *Chat) archiveMessage(message cachedMessage) error {
	valid := message.Reasoning != ""
	messageInDB, err := c.runtime.Store().DBQueries.AddMessage(c.runtime.Context(), database.AddMessageParams{
		Reasoning: sql.NullString{String: message.Reasoning, Valid: valid},
		Content:   message.Content,
		ChatID:    c.ID,
		AuthorID:  message.authorID,
		Role:      string(message.Role),
	})
	if err != nil {
		return err
	}

	message.id = messageInDB.ID

	c.cachedMessages = append(c.cachedMessages, message)

	return nil
}

func (c *Chat) archiveCurrentTurn(userMessage, assistantMessage cachedMessage) error {
	err := c.archiveMessage(userMessage)
	if err != nil {
		return err
	}

	err = c.archiveMessage(assistantMessage)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chat) getLastResponse() (cachedMessage, error) {
	var lastResponse cachedMessage
	if len(c.cachedMessages) >= 1 {
		if c.lastResponseIdx == 0 {
			lastResponse = c.cachedMessages[len(c.cachedMessages)-1]
		} else if c.lastResponseIdx > 0 {
			lastResponse = c.cachedMessages[len(c.cachedMessages)-1].variants[c.lastResponseIdx-1]
		} else {
			return cachedMessage{}, fmt.Errorf("unexpected: lastResponseIdx should be non-negative, is %d", c.lastResponseIdx)
		}
	}
	return lastResponse, nil
}

func asComMessages(messages []cachedMessage) []communication.Message {
	comMessages := []communication.Message{}
	for _, message := range messages {
		comMessages = append(comMessages, message.Message)
	}
	return comMessages
}

func (c *Chat) loadData() error {
	db := c.runtime.Store().DBQueries
	ctx := c.runtime.Context()

	messages, err := db.GetChatHistory(ctx, c.ID)
	if err != nil {
		return err
	}

	var variants []database.GetVariantsRow
	if len(messages) > 0 && messages[len(messages)-1].Role == string(communication.RoleAssistant) {
		variants, err = db.GetVariants(ctx, uuid.NullUUID{UUID: messages[len(messages)-1].ID, Valid: true})
		if err != nil {
			return err
		}
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

	var systemPrompt strings.Builder
	addToSystemPrompt := systemPromptBuilder(c.llmClient.SystemPrompt, &systemPrompt)

	fmt.Fprintf(&systemPrompt, "---\n")
	fmt.Fprintf(&systemPrompt, "[START OF CHARACTER AND CONTEXT DATA]\n")
	fmt.Fprintf(&systemPrompt, "### WORLD CONTEXT\n")
	if err = addToSystemPrompt(c.card, "chat"); err != nil {
		return err
	}

	fmt.Fprintf(&systemPrompt, "### PLAYER CHARACTER\n")
	if err = addToSystemPrompt(userCharacter.Card, "character"); err != nil {
		return err
	}

	fmt.Fprintf(&systemPrompt, "### NON-PLAYER CHARACTERS\n")
	for _, char := range c.characters {
		if err = addToSystemPrompt(char.Card, "character"); err != nil {
			return err
		}
	}

	fmt.Fprintf(&systemPrompt, "[END OF CHARACTER AND CONTEXT DATA]\n")
	fmt.Fprintf(&systemPrompt, "---\n")

	c.cachedMessages = append(c.cachedMessages, cachedMessage{
		Message: communication.Message{
			AuthorName: SystemUserName,
			Role:       communication.RoleSystem,
			Reasoning:  "",
			Content:    systemPrompt.String(),
		},
	})

	for _, message := range messages {
		role := communication.Role(message.Role)
		if !role.IsValid() {
			return fmt.Errorf("message role %q is not a valid role", role)
		}
		//fmt.Println(role)
		c.cachedMessages = append(c.cachedMessages, cachedMessage{
			id:       message.ID,
			authorID: message.AuthorID,
			Message: communication.Message{
				AuthorName: message.AuthorName,
				Role:       role,
				Reasoning:  message.Reasoning.String,
				Content:    message.Content,
			},
		})
	}

	for _, variant := range variants {
		c.cachedMessages[len(c.cachedMessages)-1].variants = append(c.cachedMessages[len(c.cachedMessages)-1].variants, cachedMessage{
			id:       variant.ID,
			authorID: variant.AuthorID,
			Message: communication.Message{
				AuthorName: variant.AuthorName,
				Role:       communication.Role(variant.Role),
				Reasoning:  variant.Reasoning.String,
				Content:    variant.Content,
			},
		})
	}

	if len(c.cachedMessages) > 1 {
		if len(variants) > 0 {
			c.lastResponseIdx = len(c.cachedMessages[len(c.cachedMessages)-1].variants)
		} else {
			c.lastResponseIdx = 0
		}
	}

	c.systemPrompt = systemPrompt.String()
	return nil
}

func systemPromptBuilder(initialSystemPrompt string, systemPromptPtr *strings.Builder) func(card json.RawMessage, cardType string) error {
	fmt.Fprintf(systemPromptPtr, "%s\n", initialSystemPrompt)

	return func(card json.RawMessage, cardType string) error {
		cardStruct := core.Card{}
		err := json.Unmarshal(card, &cardStruct)
		if err != nil {
			return err
		}

		switch cardType {
		case "character":
			if name := strings.TrimSpace(cardStruct.Name); name == "" {
				return fmt.Errorf("character with empty names are not allowed")
			}
			fmt.Fprintf(systemPromptPtr, "[Character: %s]\n", cardStruct.Name)
		case "chat":
		default:
			return fmt.Errorf("unexpected error: got unexpected card type when building system prompt")
		}

		systemPromptPtr.WriteString(formatField("Description", cardStruct.Description))
		systemPromptPtr.WriteString(formatField("Personality", cardStruct.Personality))
		systemPromptPtr.WriteString(formatField("Example Message", cardStruct.MsgExample))
		fmt.Fprintln(systemPromptPtr, "")

		//fmt.Printf("PRRRRORMPT: %s\n\n", systemPromptPtr.String())

		return nil
	}
}

func formatField(fieldName, content string) string {
	cleanedContent := strings.TrimSpace(content)
	if cleanedContent == "" {
		return ""
	}

	return fmt.Sprintf("%s: \n\"\"\"\n%s\n\"\"\"\n", fieldName, cleanedContent)
}

func mapFromDBUserCharacter(userCharacter database.GetUserCharacterInChatByIDRow) character {
	return character{
		ID:   userCharacter.ID,
		Name: userCharacter.Name,
		Card: userCharacter.Card,
	}
}
