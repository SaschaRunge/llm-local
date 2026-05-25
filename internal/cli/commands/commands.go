package commands

import (
	"github.com/SaschaRunge/llm-local/internal/core"
)

func All() []core.Command {
	return []core.Command{
		&add{},
		&atCharacter{},
		&characters{},
		&chats{},
		&debugDelete{},
		&exit{},
		&join{},
		&lobby{},
		&new{},
		&regenerate{},
	}
}
