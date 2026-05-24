package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type lobby struct{}

func (l *lobby) Name() string { return "/lobby" }
func (l *lobby) Description() string {
	return "Enters the lobby."
}
func (l *lobby) Usage() string     { return fmt.Sprintf("%s", l.Name()) }
func (l *lobby) MinArguments() int { return 0 }
func (l *lobby) MaxArguments() int { return 0 }

func (l *lobby) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (l *lobby) Execute(commandCtx core.CommandContext) (core.Result, error) {
	return core.Result{Response: "", NextScene: &scenes.Lobby{}}, nil
}
