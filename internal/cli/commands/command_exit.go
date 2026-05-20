package commands

import (
	"fmt"
	"os"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type CommandExit struct{}

func (c *CommandExit) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *CommandExit) Execute(commandCtx core.CommandContext) (core.CommandResult, error) {
	fmt.Println("=== " + core.Goodbye + " ===")
	os.Exit(0)

	return core.CommandResult{}, nil
}
