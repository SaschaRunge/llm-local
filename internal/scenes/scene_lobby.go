package scenes

import "github.com/SaschaRunge/llm-local/internal/core"

type SceneLobby struct {
}

func (c *SceneLobby) Execute(userInput string) (core.SceneResult, error) {
	return core.SceneResult{
		NextScene: c,
		Response:  "",
	}, nil
}

func (c *SceneLobby) GetName() string {
	return "Lobby"
}
