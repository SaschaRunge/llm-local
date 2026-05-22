package scenes

import "github.com/SaschaRunge/llm-local/internal/core"

type Dummy struct {
}

func (d *Dummy) Execute(userInput string) (core.Result, error) {
	return core.Result{
		NextScene: d,
		Response:  "",
	}, nil
}

func (d *Dummy) GetName() string {
	return "DummyScene"
}

func (d *Dummy) OnEnter() string {
	return ""
}
