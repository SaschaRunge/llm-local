package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type CommandChat struct{}

//should return error and replace in cli
/*
	if !exists {
		return core.CommandResult{}, fmt.Errorf("unknown command %q: %w", cmdName, core.ErrNotACommand)
	}

	if len(args) < cmdInfo.minAmountArguments {
		return core.CommandResult{}, core.ErrInvalidCommand{Context: fmt.Sprintf("not enough arguments in %q command. usage: %q", cmdInfo.name, cmdInfo.usage)}
	}
	if len(args) > cmdInfo.maxAmountArguments {
		return core.CommandResult{}, core.ErrInvalidCommand{Context: fmt.Sprintf("to many arguments in %q command. usage: %q", cmdInfo.name, cmdInfo.usage)}
	}
*/
func (c *CommandChat) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *CommandChat) Execute(commandCtx core.CommandContext) (core.CommandResult, error) {
	chats, err := commandCtx.Runtime.DB().GetChatsLikeName(commandCtx.Runtime.Context(), commandCtx.Args[0])
	if err != nil {
		return core.CommandResult{}, err
	}
	if len(chats) == 0 {
		return core.CommandResult{}, fmt.Errorf("no chat with name %q", commandCtx.Args[0])
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
		return core.CommandResult{}, err
	}

	return core.CommandResult{Output: "", NextScene: sceneChat}, nil
}
