package scenes

import "github.com/SaschaRunge/llm-local/internal/core"

type Lobby struct {
}

func (l *Lobby) Execute(userInput string) (core.CommandResult, error) {
	return core.CommandResult{
		NextScene: l,
		Output:    "",
	}, nil
}

func (l *Lobby) GetName() string {
	return "Lobby"
}
