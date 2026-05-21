package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type DebugDelete struct{}

func (d *DebugDelete) Name() string { return "/debugdelete" }
func (d *DebugDelete) Description() string {
	return "Deletes the current chat's history."
}
func (d *DebugDelete) Usage() string           { return fmt.Sprintf("%s", d.Name()) }
func (d *DebugDelete) MinAmountArguments() int { return 0 }
func (d *DebugDelete) MaxAmountArguments() int { return 0 }

func (c *DebugDelete) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (d *DebugDelete) Execute(commandCtx core.CommandContext) (core.CommandResult, error) {
	err := commandCtx.Runtime.DB().DeleteMessages(commandCtx.Runtime.Context())
	if err != nil {
		return core.CommandResult{}, err
	}

	fmt.Println("Debug Delete Messages.")
	return core.CommandResult{}, nil
}
