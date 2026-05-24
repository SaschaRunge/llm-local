package commands

import (
	"fmt"
	"os"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type exit struct{}

func (e *exit) Name() string { return "/exit" }
func (e *exit) Description() string {
	return "Exits the program."
}
func (e *exit) Usage() string     { return fmt.Sprintf("%s", e.Name()) }
func (e *exit) MinArguments() int { return 0 }
func (e *exit) MaxArguments() int { return 0 }

func (e *exit) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (e *exit) Execute(commandCtx core.CommandContext) (core.Result, error) {
	fmt.Println("=== " + core.Goodbye + " ===")
	os.Exit(0)

	return core.Result{}, nil
}
