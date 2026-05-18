package cli

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

func executeCommandChat(ctx commandContext) (core.Scene, error) {
	chats, err := ctx.cli.DBQueries.GetChatsLikeName(ctx.cli.context, ctx.args[0])
	if err != nil {
		return nil, err
	}
	if len(chats) == 0 {
		return nil, fmt.Errorf("no chat with name %q", ctx.args[0])
	}

	/*// TODO: selection from multiple personas
	personas, err := ctx.cli.DBQueries.GetPersonasLikeName(ctx.cli.context, ctx.args[1])
	if err != nil {
		return nil, err
	}
	if len(personas) == 0 {
		return nil, fmt.Errorf("no persona with name %q", ctx.args[1])
	}*/

	sceneChat, err := scenes.NewSceneChat(ctx.cli, chats[0]) /*, scenes.Character{
		ID:   personas[0].ID,
		Name: personas[0].Name,
	})*/
	if err != nil {
		return nil, err
	}

	return sceneChat, nil
}
