package commands

import (
	"github.com/SaschaRunge/llm-local/internal/core"
)

func All() []core.Command {
	return []core.Command{
		&atCharacter{},
		&characters{},
		&join{},
		&chats{},
		&debugDelete{},
		&exit{},
		&lobby{},
		&new{},
		&regenerate{},
	}
}
