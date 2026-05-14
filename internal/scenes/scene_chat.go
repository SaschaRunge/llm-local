package scenes

import (
	"github.com/SaschaRunge/llm-local/internal/core"

	"github.com/google/uuid"
)

type SceneChat struct {
	messages   []message
	characters []character
	ID         uuid.UUID
	Name       string
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

func (c *SceneChat) Run(runtime core.Runtime) (core.Scene, error) {
	err := c.loadData(runtime)
	if err != nil {
		return nil, err
	}

	for {
		userInput := runtime.GetInput()
		next, err := runtime.ExecuteCommand(userInput)
		if next != nil || err != nil {
			return next, err
		}
	}
}

func (c *SceneChat) loadData(runtime core.Runtime) error {
	db := runtime.DB()

	messages, err := db.GetChatHistory(runtime.Context(), c.ID)
	if err != nil {
		return err
	}

	characters, err := db.GetCharactersInChatByChatID(runtime.Context(), c.ID)
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
