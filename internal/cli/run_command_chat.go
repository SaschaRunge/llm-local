package cli

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/scenes"
)

func runCommandChat(ctx commandContext) error {
	chats, err := ctx.cli.DBQueries.GetChatsLikeName(ctx.cli.context, ctx.args[0])
	if err != nil {
		return nil
	}
	if len(chats) == 0 {
		return fmt.Errorf("No chat with name '%s'.", ctx.args[0])
	}
	sceneChat := scenes.SceneChat{
		ID:   chats[0].ID,
		Name: chats[0].Name,
	}
	sceneChat.Run(ctx.cli)

	return nil
}
