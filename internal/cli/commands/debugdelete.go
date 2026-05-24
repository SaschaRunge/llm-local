package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type debugDelete struct{}

func (d *debugDelete) Name() string { return "/debugdelete" }
func (d *debugDelete) Description() string {
	return "Deletes the current chat's history."
}
func (d *debugDelete) Usage() string     { return fmt.Sprintf("%s", d.Name()) }
func (d *debugDelete) MinArguments() int { return 0 }
func (d *debugDelete) MaxArguments() int { return 0 }

func (c *debugDelete) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (d *debugDelete) Execute(commandCtx core.CommandContext) (core.Result, error) {
	err := commandCtx.Runtime.DB().DeleteMessages(commandCtx.Runtime.Context())
	if err != nil {
		return core.Result{}, err
	}

	fmt.Println("Debug Delete Messages.")
	return core.Result{}, nil
}
