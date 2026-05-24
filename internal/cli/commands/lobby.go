package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type Lobby struct{}

func (r *Lobby) Name() string { return "/lobby" }
func (r *Lobby) Description() string {
	return "Enters the lobby."
}
func (r *Lobby) Usage() string     { return fmt.Sprintf("%s", r.Name()) }
func (r *Lobby) MinArguments() int { return 0 }
func (r *Lobby) MaxArguments() int { return 0 }

func (r *Lobby) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (r *Lobby) Execute(commandCtx core.CommandContext) (core.Result, error) {
	return core.Result{Response: "", NextScene: &scenes.Lobby{}}, nil
}
