package scenes

import "github.com/SaschaRunge/llm-local/internal/core"

type SceneDummy struct {
}

func (c *SceneDummy) Execute(userInput string) (core.SceneResult, error) {
	return core.SceneResult{
		NextScene: c,
		Response:  "",
	}, nil
}

func (c *SceneDummy) GetName() string {
	return "DummyScene"
}
