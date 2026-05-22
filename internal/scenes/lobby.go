package scenes

import "github.com/SaschaRunge/llm-local/internal/core"

type Lobby struct {
}

func (l *Lobby) Execute(userInput string) (core.Result, error) {
	return core.Result{
		NextScene: l,
		Response:  "",
	}, nil
}

func (l *Lobby) GetName() string {
	return "Lobby"
}

func (l *Lobby) OnEnter() string {
	return ""
}
