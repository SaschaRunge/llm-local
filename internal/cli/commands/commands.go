package commands

import (
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
)

type option func(core.CommandContext) (core.Result, error)

type character = database.GetAvailableCharactersRow
type chat = database.GetAvailableChatsRow

func All() []core.Command {
	return []core.Command{
		&add{},
		&atCharacter{},
		&characters{},
		&chats{},
		&cycle{},
		&debugDelete{},
		&edit{},
		&exit{},
		&join{},
		&lobby{},
		&new{},
		&regenerate{},
	}
}
