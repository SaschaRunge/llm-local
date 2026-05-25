package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type lobby struct{}

func (c *lobby) Name() string { return "/lobby" }
func (c *lobby) Description() string {
	return "Enters the lobby."
}
func (c *lobby) Usage() string     { return fmt.Sprintf("%s", c.Name()) }
func (c *lobby) MinArguments() int { return 0 }
func (c *lobby) MaxArguments() int { return 0 }

func (c *lobby) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *lobby) Execute(commandCtx core.CommandContext) (core.Result, error) {
	return core.Result{Response: "", NextScene: &scenes.Lobby{}}, nil
}
