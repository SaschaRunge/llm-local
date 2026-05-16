package cli

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

func runCommandChat(ctx commandContext) (core.Scene, error) {
	chats, err := ctx.cli.DBQueries.GetChatsLikeName(ctx.cli.context, ctx.args[0])
	if err != nil {
		return nil, err
	}
	if len(chats) == 0 {
		return nil, fmt.Errorf("No chat with name '%s'.", ctx.args[0])
	}
	sceneChat, err := scenes.NewSceneChat(ctx.cli, chats[0])
	if err != nil {
		return nil, err
	}

	return &sceneChat, nil
}
