package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
)

type CommandChats struct{}

func (c *CommandChats) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *CommandChats) Execute(commandCtx core.CommandContext) (core.CommandResult, error) {
	chats, err := commandCtx.Runtime.DB().GetAllChats(commandCtx.Runtime.Context())
	if err != nil {
		return core.CommandResult{Output: "", NextScene: commandCtx.Runtime.CurrentScene()}, err
	}

	if len(chats) == 0 {
		fmt.Println("No available chats.")
		return core.CommandResult{Output: "", NextScene: commandCtx.Runtime.CurrentScene()}, nil
	}

	fmt.Println("Available chats:")
	for i, chat := range chats {
		fmt.Printf("%d. %s\n", i, chat.Name)
	}

	return core.CommandResult{Output: "", NextScene: commandCtx.Runtime.CurrentScene()}, nil
}
