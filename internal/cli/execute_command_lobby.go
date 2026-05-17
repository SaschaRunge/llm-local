package cli

import (
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

func executeCommandLobby(ctx commandContext) (core.Scene, error) {
	return &scenes.SceneLobby{}, nil
}
