package commands

import (
	"fmt"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type characters struct{}

func (c *characters) Name() string { return "/characters" }
func (c *characters) Description() string {
	return "Shows the available characters in the current chat."
}
func (c *characters) Usage() string     { return fmt.Sprintf("%s", c.Name()) }
func (c *characters) MinArguments() int { return 0 }
func (c *characters) MaxArguments() int { return 0 }

func (c *characters) CanExecute(commandCtx core.CommandContext) bool {
	_, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	return isSceneChat
}

func (c *characters) Execute(commandCtx core.CommandContext) (core.Result, error) {
	chat, ok := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	if !ok {
		return core.Result{}, fmt.Errorf("unexpected error in command %q: type assertion failed", c.Name())
	}

	availableCharacters := chat.GetAvailableCharacters()
	if len(availableCharacters) == 0 {
		return core.Result{}, fmt.Errorf("no characters in chat %q", c.Name())
	}

	var response strings.Builder

	for i, char := range availableCharacters {
		_, err := fmt.Fprintf(&response, "%d. %s\n", i+1, char)
		if err != nil {
			return core.Result{}, fmt.Errorf("unexpected error in command %q: %w", c.Name(), err)
		}
	}

	return core.Result{Response: fmt.Sprintf("The following characters are currently in chat %s:\n%s", chat.GetName(), response.String())}, nil
}
