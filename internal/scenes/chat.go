package scenes

import (
	"context"
	"database/sql"
	"fmt"
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

type Chat struct {
	cachedMessages []cachedMessage
	characters     []Character
	userCharacter  Character
	ID             uuid.UUID
	Name           string
	scenario       string

	runtime   core.Runtime
	llmClient *llm.Client
}

type cachedMessage struct {
	id       uuid.UUID
	authorID uuid.UUID
	communication.Message
}

type Character = database.GetCharactersInChatRow

// TODO: move to fallback
type CommandExecuteChat struct{}

func (c *CommandExecuteChat) CanExecute(commandCtx core.CommandContext) bool {
	_, ok := commandCtx.Runtime.CurrentScene().(*Chat)
	return ok
}

func (c *CommandExecuteChat) Execute(commandCtx core.CommandContext) (core.Result, error) {
	return commandCtx.Runtime.CurrentScene().Execute(commandCtx.Cmd)
}

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

// TODO: need to add id/authorID to cached Messages after commiting to DB
// instead append to history at the end
func (c *Chat) Execute(userInput string) (core.Result, error) {
	history := append([]cachedMessage{}, c.cachedMessages...)
	prompt := cachedMessage{
		Message: communication.Message{
			Name:      c.userCharacter.Name,
			Role:      communication.RoleUser,
			Reasoning: "",
			Content:   userInput,
		},
	}

	history = append(history, prompt)

	ctx, cancel := context.WithTimeout(c.runtime.Context(), generateAnswerTimeoutInSeconds*time.Second)
	defer cancel()
	response, err := c.llmClient.GenerateAnswer(ctx, asComMessages(history))
	if err != nil {
		return core.Result{
			Response:  "",
			NextScene: c,
		}, err
	}

	answer := cachedMessage{
		Message: communication.Message{
			Name:      c.characters[0].Name,
			Role:      communication.RoleAssistant,
			Reasoning: "",
			Content:   response,
		},
	}

	//TODO: rollback entire commit if send fails on either Add
	promptInDB, err := c.runtime.DB().AddMessage(c.runtime.Context(), database.AddMessageParams{
		Reasoning: sql.NullString{},
		Content:   prompt.Content,
		ChatID:    c.ID,
		AuthorID:  c.userCharacter.ID,
		Role:      string(prompt.Role),
	})
	if err != nil {
		return core.Result{
			Response:  "",
			NextScene: c,
		}, err
	}

	//TODO: currently answers just as the first character in the character list
	answerInDB, err := c.runtime.DB().AddMessage(c.runtime.Context(), database.AddMessageParams{
		Reasoning: sql.NullString{},
		Content:   answer.Content,
		ChatID:    c.ID,
		AuthorID:  c.characters[0].ID,
		Role:      string(answer.Role),
	})
	if err != nil {
		return core.Result{
			Response:  "",
			NextScene: c,
		}, err
	}

	prompt.id = promptInDB.ID
	prompt.authorID = promptInDB.AuthorID
	answer.id = answerInDB.ID
	answer.authorID = answerInDB.AuthorID

	c.cachedMessages = append(c.cachedMessages, prompt, answer)

	return core.Result{
		Response:  response,
		NextScene: c,
	}, nil
}

func (c *Chat) GetName() string {
	return "Chat"
}

func (c *Chat) Regenerate() core.Result {
	return core.Result{Response: "Totally regenerated!", NextScene: c}
}

func (c *Chat) loadData() error {
	db := c.runtime.DB()
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

	scenario, err := db.GetScenarioFromChatByID(ctx, c.ID)
	if err != nil {
		return err
	}
	c.scenario = scenario.String

	for _, message := range messages {
		role := communication.Role(message.Role)
		if !role.IsValid() {
			return fmt.Errorf("message role %q is not a valid role", role)
		}
		cachedMessage := cachedMessage{
			id:       message.ID,
			authorID: message.AuthorID,
			Message: communication.Message{
				Name:      message.AuthorName,
				Role:      role,
				Reasoning: message.Reasoning.String,
				Content:   message.Content,
			},
		}
		c.cachedMessages = append(c.cachedMessages, cachedMessage)
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

func mapFromDBUserCharacter(userCharacter database.GetUserCharacterInChatByIDRow) Character {
	return Character{
		ID:           userCharacter.ID,
		Name:         userCharacter.Name,
		SystemPrompt: userCharacter.SystemPrompt,
	}
}

func (c *Chat) writeToHistory() error {
	return nil
}
