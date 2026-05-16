package scenes

import (
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"

	"github.com/google/uuid"
)

type SceneChat struct {
	messages   []message
	characters []character
	ID         uuid.UUID
	Name       string
	runtime    core.Runtime
}

type message struct {
	id              uuid.UUID
	authorID        uuid.UUID
	authorName      string
	contentThoughts string
	contentAnswer   string
}

type character struct {
	id   uuid.UUID
	name string
}

func NewSceneChat(runtime core.Runtime, chat database.Chat) (SceneChat, error) {
	sceneChat := SceneChat{
		runtime: runtime,
		ID:      chat.ID,
		Name:    chat.Name}
	if err := sceneChat.loadData(); err != nil {
		return SceneChat{}, err
	}

	return sceneChat, nil
}

// TODO: get rid of handle
func (c *SceneChat) Execute(userInput string) (core.Scene, error) {
	next, err := c.runtime.Handle(userInput)
	if next != nil || err != nil {
		return next, err
	}

	//TODO: llm logic

	return c, nil
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
		c.messages = append(c.messages, message{
			id:              msg.ID,
			authorID:        msg.AuthorID,
			authorName:      msg.AuthorName,
			contentThoughts: msg.ContentThoughts.String,
			contentAnswer:   msg.ContentAnswer,
		})
	}

	return nil
}
