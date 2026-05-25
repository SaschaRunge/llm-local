package commands

import (
	"fmt"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/cli/formatters"
	"github.com/SaschaRunge/llm-local/internal/core"
)

type chats struct{}

func (c *chats) Name() string { return "/chats" }
func (c *chats) Description() string {
	return "Shows the available chats for the current user."
}
func (c *chats) Usage() string     { return fmt.Sprintf("%s", c.Name()) }
func (c *chats) MinArguments() int { return 0 }
func (c *chats) MaxArguments() int { return 0 }

func (c *chats) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *chats) Execute(commandCtx core.CommandContext) (core.Result, error) {
	chats, err := commandCtx.Runtime.DB().GetAllChats(commandCtx.Runtime.Context())
	if err != nil {
		return core.Result{}, err
	}

	if len(chats) == 0 {
		fmt.Println("No available chats.")
		return core.Result{}, nil
	}

	/*
		var response strings.Builder

		for i, chat := range chats {
			_, err := fmt.Fprintf(&response, "  %d. %s\n", i+1, chat.Name)
			if err != nil {
				return core.Result{}, fmt.Errorf("unexpected error in command %q: %w", c.Name(), err)
			}
		}
	*/

	var response strings.Builder
	addListItem := formatters.ListBuilder(&response)

	for _, chat := range chats {
		addListItem(chat.Name)
	}

	return core.Result{Response: response.String()}, nil
}
