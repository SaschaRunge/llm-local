package modes

import (
	"github.com/google/uuid"
)

type ModeChat struct {
	characters []struct {
		id   uuid.UUID
		name string
	}
	messages []struct {
		id               uuid.UUID
		author_id        uuid.UUID
		author_name      string
		content_thoughts string
		content_answer   string
	}
}

func (c *ModeChat) Run() {
	for {

	}
}
