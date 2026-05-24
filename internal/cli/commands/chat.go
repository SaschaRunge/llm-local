package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type Chat struct{}

func (c *Chat) Name() string { return "/chat" }
func (c *Chat) Description() string {
	return "Starts the chat interface. Will load or create the chat [chat_name]."
}
func (c *Chat) Usage() string     { return fmt.Sprintf("%s [chat_name]", c.Name()) }
func (c *Chat) MinArguments() int { return 1 }
func (c *Chat) MaxArguments() int { return 1 }

func (c *Chat) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *Chat) Execute(commandCtx core.CommandContext) (core.Result, error) {
	chats, err := commandCtx.Runtime.DB().GetChatsLikeName(commandCtx.Runtime.Context(), commandCtx.Args[0])
	if err != nil {
		return core.Result{}, err
	}
	if len(chats) == 0 {
		return core.Result{}, fmt.Errorf("no chat with name %q", commandCtx.Args[0])
	}

	/*// TODO: selection from multiple personas
	personas, err := ctx.cli.DBQueries.GetPersonasLikeName(ctx.cli.context, ctx.args[1])
	if err != nil {
		return nil, err
	}
	if len(personas) == 0 {
		return nil, fmt.Errorf("no persona with name %q", ctx.args[1])
	}*/

	sceneChat, err := scenes.NewSceneChat(commandCtx.Runtime, chats[0]) /*, scenes.Character{
		ID:   personas[0].ID,
		Name: personas[0].Name,
	})*/
	if err != nil {
		return core.Result{}, err
	}

	return core.Result{Response: "", NextScene: sceneChat}, nil
}
