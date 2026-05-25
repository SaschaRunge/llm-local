package commands

import (
	"fmt"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type atCharacter struct{}

func (c *atCharacter) Name() string { return "@" }
func (c *atCharacter) Description() string {
	return "Generates an answer as the specified character."
}
func (c *atCharacter) Usage() string     { return "@[character] [text]" }
func (c *atCharacter) MinArguments() int { return 2 }
func (c *atCharacter) MaxArguments() int { return 2 }

func (c *atCharacter) CanExecute(commandCtx core.CommandContext) bool {
	_, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	return isSceneChat
}

func (c *atCharacter) Execute(commandCtx core.CommandContext) (core.Result, error) {
	chat, ok := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	if !ok {
		return core.Result{}, fmt.Errorf("unexpected error in command @character: type assertion failed")
	}

	result, err := chat.AtCharacter(commandCtx.Args[0], commandCtx.Args[1])
	if err != nil {
		err = fmt.Errorf("mention character %q failed with error: %w", commandCtx.Args[0], err)
	}

	return result, err
}

// ParseArgs returns charactername at index 0, prompt at index 1
func (c *atCharacter) ParseArgs(rawArgs string) []string {
	parts := strings.SplitN(rawArgs, " ", 2)
	return parts
}
