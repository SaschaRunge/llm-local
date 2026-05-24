package commands

import (
	"github.com/SaschaRunge/llm-local/internal/core"
)

func All() []core.Command {
	return []core.Command{
		&atCharacter{},
		&characters{},
		&chat{},
		&chats{},
		&debugDelete{},
		&exit{},
		&lobby{},
		&regenerate{},
	}
}
