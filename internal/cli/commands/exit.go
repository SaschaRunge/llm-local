package commands

import (
	"fmt"
	"os"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type Exit struct{}

func (l *Exit) Name() string { return "/exit" }
func (l *Exit) Description() string {
	return "Exits the program."
}
func (l *Exit) Usage() string           { return fmt.Sprintf("%s", l.Name()) }
func (l *Exit) MinAmountArguments() int { return 0 }
func (l *Exit) MaxAmountArguments() int { return 0 }

func (l *Exit) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (l *Exit) Execute(commandCtx core.CommandContext) (core.CommandResult, error) {
	fmt.Println("=== " + core.Goodbye + " ===")
	os.Exit(0)

	return core.CommandResult{}, nil
}
