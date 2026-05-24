package commands

import (
	"github.com/SaschaRunge/llm-local/internal/core"
)

func All() []core.Command {
	return []core.Command{
		&atCharacter{},
		&Chat{},
		&Chats{},
		&DebugDelete{},
		&Exit{},
		&Lobby{},
		&Regenerate{},
	}
}
