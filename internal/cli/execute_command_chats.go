package cli

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
)

func executeCommandChats(ctx commandContext) (core.Scene, error) {
	chats, err := ctx.cli.DBQueries.GetAllChats(ctx.cli.context)
	if err != nil {
		return ctx.cli.CurrentScene, err
	}

	if len(chats) == 0 {
		fmt.Println("No available chats.")
		return ctx.cli.CurrentScene, nil
	}

	fmt.Println("Available chats:")
	for i, chat := range chats {
		fmt.Printf("%d. %s\n", i, chat.Name)
	}

	return ctx.cli.CurrentScene, nil
}
