package scenes

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/communication"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/SaschaRunge/llm-local/internal/llm"

	"github.com/google/uuid"
)

type SceneChat struct {
	cachedMessages []cachedMessage
	characters     []character
	ID             uuid.UUID
	Name           string
	runtime        core.Runtime
	llmClient      *llm.Client
}

type cachedMessage struct {
	id       uuid.UUID
	authorID uuid.UUID
	communication.Message
}

type character struct {
	id   uuid.UUID
	name string
}

func NewSceneChat(runtime core.Runtime, chat database.Chat) (*SceneChat, error) {
	sceneChat := SceneChat{
		runtime: runtime,
		ID:      chat.ID,
		Name:    chat.Name,
	}

	if err := sceneChat.loadData(); err != nil {
		return &SceneChat{}, err
	}

	//TODO: do not hardcode model type
	llmClient, err := llm.NewClient("qwen3.6-27b")
	sceneChat.llmClient = llmClient

	if err != nil {
		return &SceneChat{}, err
	}

	return &sceneChat, nil
}

// TODO: get rid of handle
func (c *SceneChat) Execute(userInput string) (core.SceneResult, error) {

	return core.SceneResult{
		Response:  "",
		NextScene: c,
	}, nil
}

func (c *SceneChat) GetName() string {
	return "Chat"
}

func (c *SceneChat) loadData() error {
	db := c.runtime.DB()

	messages, err := db.GetChatHistory(c.runtime.Context(), c.ID)
	if err != nil {
		return err
	}

	characters, err := db.GetCharactersInChat(c.runtime.Context(), c.ID)
	if err != nil {
		return err
	}

	for _, char := range characters {
		c.characters = append(c.characters, character{
			id:   char.ID,
			name: char.Name,
		})
	}

	for _, msg := range messages {
		role := communication.Role(msg.Role)
		if !role.IsValid() {
			return fmt.Errorf("message role %q is not a valid role", role)
		}
		cachedMessage := cachedMessage{
			id:       msg.ID,
			authorID: msg.AuthorID,
			Message: communication.Message{
				Name:      msg.AuthorName,
				Role:      role,
				Reasoning: msg.ContentThoughts.String,
				Content:   msg.ContentAnswer,
			},
		}
		c.cachedMessages = append(c.cachedMessages, cachedMessage)
	}

	return nil
}

func (c *SceneChat) writeToHistory() error {
	return nil
}
