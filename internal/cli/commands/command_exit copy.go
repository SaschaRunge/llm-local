package commands

import (
	"github.com/SaschaRunge/llm-local/internal/core"
)

type CommandLobby struct{}

func (c *CommandLobby) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *CommandLobby) Execute(commandCtx core.CommandContext) (core.CommandResult, error) {
	return core.CommandResult{}, nil
}
