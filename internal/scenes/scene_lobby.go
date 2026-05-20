package scenes

import "github.com/SaschaRunge/llm-local/internal/core"

type SceneLobby struct {
}

func (c *SceneLobby) Execute(userInput string) (core.CommandResult, error) {
	return core.CommandResult{
		NextScene: c,
		Output:    "",
	}, nil
}

func (c *SceneLobby) GetName() string {
	return "Lobby"
}
