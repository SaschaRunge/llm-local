package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type regenerate struct{}

func (r *regenerate) Name() string { return "/regenerate" }
func (r *regenerate) Description() string {
	return "Regenerates message [id]."
}
func (r *regenerate) Usage() string     { return fmt.Sprintf("%s", r.Name()) }
func (r *regenerate) MinArguments() int { return 1 }
func (r *regenerate) MaxArguments() int { return 1 }

func (r *regenerate) CanExecute(commandCtx core.CommandContext) bool {
	_, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	return isSceneChat
}

func (r *regenerate) Execute(commandCtx core.CommandContext) (core.Result, error) {
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

func (r *regenerate) ParseArgs(rawArgs string) []string {
	return []string{rawArgs}
}
