package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/cli/commands/selector"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
)

func selectChat(commandCtx core.CommandContext, dbQueries *database.Queries) (chat, error) {
	chats, err := dbQueries.GetAvailableChats(commandCtx.Runtime.Context())
	if err != nil {
		return chat{}, err
	}
	if len(chats) == 0 {
		return chat{}, fmt.Errorf("no valid chats available")
	}

	selectedChat, err := selector.Select(
		chats,
		func(chat chat) string { return chat.Name },
		commandCtx.Runtime.GetInput)
	if err != nil {
		return chat{}, err
	}
	if selectedChat.Name == "" {
		return chat{}, fmt.Errorf("chat not found")
	}

	return selectedChat, nil
}
