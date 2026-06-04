package commands

import (
	"fmt"
	"strconv"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type cycle struct{}

func (c *cycle) Name() string { return "/cycle" }
func (c *cycle) Description() string {
	return "Cycle through available regenerated variants."
}
func (c *cycle) Usage() string     { return fmt.Sprintf("%s [OPTIONAL: index]", c.Name()) }
func (c *cycle) MinArguments() int { return 0 }
func (c *cycle) MaxArguments() int { return 1 }

func (c *cycle) CanExecute(commandCtx core.CommandContext) bool {
	_, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	return isSceneChat
}

func (c *cycle) Execute(commandCtx core.CommandContext) (core.Result, error) {
	chat, ok := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	if !ok {
		return core.Result{}, fmt.Errorf("unexpected error in command /regenerate: type assertion failed")
	}

	if len(commandCtx.Args) == 0 {
		return chat.Cycle(-1)
	} else {
		idx, err := strconv.Atoi(commandCtx.Args[0])
		if err != nil {
			return core.Result{}, fmt.Errorf("argument [index] must be a number")
		}
		return chat.Cycle(idx)
	}
}
