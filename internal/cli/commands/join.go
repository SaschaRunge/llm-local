package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/cli/commands/selector"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type join struct{}

func (c *join) Name() string { return "/join" }
func (c *join) Description() string {
	return "Starts the chat interface. Will load or create the chat [chat_name]."
}
func (c *join) Usage() string     { return fmt.Sprintf("%s [OPTIONAL:chat_name]", c.Name()) }
func (c *join) MinArguments() int { return 0 }
func (c *join) MaxArguments() int { return 1 }

func (c *join) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *join) Execute(commandCtx core.CommandContext) (core.Result, error) {
	var chats []database.Chat
	var selectedChat database.Chat
	var err error
	if len(commandCtx.Args) == 1 {
		chats, err = commandCtx.Runtime.Store().DBQueries.GetChatsLikeName(commandCtx.Runtime.Context(), commandCtx.Args[0])
		if err != nil {
			return core.Result{}, err
		}
		if len(chats) == 0 {
			return core.Result{}, fmt.Errorf("no chat with name %q", commandCtx.Args[0])
		}
		if len(chats) > 1 {
			return core.Result{}, fmt.Errorf("multiple results for %q, please try again", commandCtx.Args[0])
		}
		selectedChat = chats[0]
	} else {
		chats, err := commandCtx.Runtime.Store().DBQueries.GetAllChats(commandCtx.Runtime.Context())
		if err != nil {
			return core.Result{}, err
		}
		if len(chats) == 0 {
			return core.Result{}, fmt.Errorf("no valid chats available")
		}

		selectedChat = selector.Select(
			chats,
			func(chat database.Chat) string { return chat.Name },
			commandCtx.Runtime.GetInput)

		if selectedChat.Name == "" {
			return core.Result{}, err
		}
	}

	sceneChat, err := scenes.NewSceneChat(commandCtx.Runtime, selectedChat)
	if err != nil {
		return core.Result{}, err
	}

	return core.Result{Response: "", NextScene: sceneChat}, nil
}
