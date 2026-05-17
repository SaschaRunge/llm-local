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
	messages   []chatMessage
	characters []character
	ID         uuid.UUID
	Name       string
	runtime    core.Runtime
	llmClient  *llm.Client
}

type chatMessage struct {
	id         uuid.UUID
	authorID   uuid.UUID
	authorName string
	reasoning  string
	content    string
	role       string
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

	messageHistory := []communication.Message{}
	for _, message := range sceneChat.messages {
		newMessage, err := translateToCommunicationMessage(message)
		if err != nil {
			return &SceneChat{}, err
		}
		messageHistory = append(messageHistory, newMessage)
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
		c.messages = append(c.messages, chatMessage{
			id:         msg.ID,
			authorID:   msg.AuthorID,
			authorName: msg.AuthorName,
			reasoning:  msg.ContentThoughts.String,
			content:    msg.ContentAnswer,
			role:       msg.Role,
		})
	}

	return nil
}

func translateToChatMessage(msg communication.Message) chatMessage {
	return chatMessage{
		authorName: msg.Name,
		reasoning:  msg.Reasoning,
		content:    msg.Content,
	}
}

func translateToCommunicationMessage(msg chatMessage) (communication.Message, error) {
	role := communication.Role(msg.role)
	if !role.IsValid() {
		return communication.Message{}, fmt.Errorf("message role %q is not a valid role", msg.role)
	}

	return communication.Message{
		Name:      msg.authorName,
		Reasoning: msg.reasoning,
		Content:   msg.content,
		Role:      role,
	}, nil
}
