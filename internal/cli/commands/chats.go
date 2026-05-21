package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type Chats struct{}

func (c *Chats) Name() string { return "/chats" }
func (c *Chats) Description() string {
	return "Shows the available chats for the current user."
}
func (c *Chats) Usage() string           { return fmt.Sprintf("%s", c.Name()) }
func (c *Chats) MinAmountArguments() int { return 0 }
func (c *Chats) MaxAmountArguments() int { return 0 }

func (c *Chats) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *Chats) Execute(commandCtx core.CommandContext) (core.CommandResult, error) {
	chats, err := commandCtx.Runtime.DB().GetAllChats(commandCtx.Runtime.Context())
	if err != nil {
		return core.CommandResult{Output: "", NextScene: commandCtx.Runtime.CurrentScene()}, err
	}

	if len(chats) == 0 {
		fmt.Println("No available chats.")
		return core.CommandResult{Output: "", NextScene: commandCtx.Runtime.CurrentScene()}, nil
	}

	// TODO: return as SPrintf

	fmt.Println("Available chats:")
	for i, chat := range chats {
		fmt.Printf("%d. %s\n", i, chat.Name)
	}

	return core.CommandResult{Output: "", NextScene: commandCtx.Runtime.CurrentScene()}, nil
}
