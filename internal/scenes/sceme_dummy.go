package scenes

import "github.com/SaschaRunge/llm-local/internal/core"

type SceneDummy struct {
}

func (c *SceneDummy) Execute(userInput string) (core.CommandResult, error) {
	return core.CommandResult{
		NextScene: c,
		Output:    "",
	}, nil
}

func (c *SceneDummy) GetName() string {
	return "DummyScene"
}
