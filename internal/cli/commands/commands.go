package commands

import (
	"github.com/SaschaRunge/llm-local/internal/core"
)

type Library struct{}

func (l *Library) GetAll() []core.Command {
	return []core.Command{
		&Chat{},
		&Chats{},
		&DebugDelete{},
		&Exit{},
		&Lobby{},
		&Regenerate{},
	}
}
