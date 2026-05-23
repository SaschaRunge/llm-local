package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type atCharacter struct{}

func (r *atCharacter) Name() string { return "@character" }
func (r *atCharacter) Description() string {
	return "Generates an answer as the specified character."
}
func (r *atCharacter) Usage() string           { return "@[Character]" }
func (r *atCharacter) MinAmountArguments() int { return 0 }
func (r *atCharacter) MaxAmountArguments() int { return 0 }

func (r *atCharacter) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (r *atCharacter) Execute(commandCtx core.CommandContext) (core.Result, error) {
	chat, ok := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	if !ok {
		return core.Result{}, fmt.Errorf("unexpected error in command @character: type assertion failed")
	}
	return chat.AtCharacter(commandCtx.Args[0], commandCtx.Args[1])
}

func (r *atCharacter) ParseArgs(rawArgs string) []string {
	return []string{rawArgs}
}
