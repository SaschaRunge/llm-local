package commands

import (
	"fmt"
	"os"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type exit struct{}

func (c *exit) Name() string { return "/exit" }
func (c *exit) Description() string {
	return "Exits the program."
}
func (c *exit) Usage() string     { return fmt.Sprintf("%s", c.Name()) }
func (c *exit) MinArguments() int { return 0 }
func (c *exit) MaxArguments() int { return 0 }

func (c *exit) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *exit) Execute(commandCtx core.CommandContext) (core.Result, error) {
	fmt.Println("=== " + core.Goodbye + " ===")
	os.Exit(0)

	return core.Result{}, nil
}
