package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type regenerate struct{}

func (c *regenerate) Name() string { return "/regenerate" }
func (c *regenerate) Description() string {
	return "Regenerates message [id]."
}
func (c *regenerate) Usage() string     { return fmt.Sprintf("%s", c.Name()) }
func (c *regenerate) MinArguments() int { return 1 }
func (c *regenerate) MaxArguments() int { return 1 }

func (c *regenerate) CanExecute(commandCtx core.CommandContext) bool {
	_, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	return isSceneChat
}

func (c *regenerate) Execute(commandCtx core.CommandContext) (core.Result, error) {
	chat, ok := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	if !ok {
		return core.Result{}, fmt.Errorf("unexpected error in command /regenerate: type assertion failed")
	}

	result, err := chat.Regenerate(commandCtx.Args[0])
	if err != nil {
		err = fmt.Errorf("regeneration failed with error: %w", err)
	}

	return result, err
}

func (c *regenerate) ParseArgs(rawArgs string) []string {
	return []string{rawArgs}
}
