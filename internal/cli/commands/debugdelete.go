package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type debugDelete struct{}

func (c *debugDelete) Name() string { return "/debugdelete" }
func (c *debugDelete) Description() string {
	return "Deletes the current chat's history."
}
func (c *debugDelete) Usage() string     { return fmt.Sprintf("%s", c.Name()) }
func (c *debugDelete) MinArguments() int { return 0 }
func (c *debugDelete) MaxArguments() int { return 0 }

func (c *debugDelete) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *debugDelete) Execute(commandCtx core.CommandContext) (core.Result, error) {
	err := commandCtx.Runtime.Store().DBQueries.DeleteMessages(commandCtx.Runtime.Context())
	if err != nil {
		return core.Result{}, err
	}

	fmt.Println("Debug Delete Messages.")
	return core.Result{}, nil
}
