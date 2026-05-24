package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type Regenerate struct{}

func (r *Regenerate) Name() string { return "/regenerate" }
func (r *Regenerate) Description() string {
	return "Regenerates message [id]."
}
func (r *Regenerate) Usage() string     { return fmt.Sprintf("%s", r.Name()) }
func (r *Regenerate) MinArguments() int { return 1 }
func (r *Regenerate) MaxArguments() int { return 1 }

func (r *Regenerate) CanExecute(commandCtx core.CommandContext) bool {
	_, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	return isSceneChat
}

func (r *Regenerate) Execute(commandCtx core.CommandContext) (core.Result, error) {
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

func (r *Regenerate) ParseArgs(rawArgs string) []string {
	return []string{rawArgs}
}
