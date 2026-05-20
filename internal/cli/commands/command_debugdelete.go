package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type CommandDebugDelete struct{}

func (c *CommandDebugDelete) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *CommandDebugDelete) Execute(commandCtx core.CommandContext) (core.CommandResult, error) {
	err := commandCtx.Runtime.DB().DeleteMessages(commandCtx.Runtime.Context())
	if err != nil {
		return core.CommandResult{}, err
	}

	fmt.Println("Debug Delete Messages.")
	return core.CommandResult{}, nil
}
