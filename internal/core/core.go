package core

import (
	"context"

	"github.com/SaschaRunge/llm-local/internal/database"
)

type CommandResult struct {
	Output    string
	NextScene Scene
}

type CommandContext struct {
	Cmd     string
	Args    []string
	Runtime Runtime
}

type Runtime interface {
	Context() context.Context
	CurrentScene() Scene
	DB() *database.Queries
}

type Scene interface {
	GetName() string
	Execute(rawInput string) (CommandResult, error)
}

type Command interface {
	Name() string
	Description() string
	Usage() string
	MinAmountArguments() int
	MaxAmountArguments() int

	CanExecute(commandCtx CommandContext) bool
	Execute(commandCtx CommandContext) (CommandResult, error)
}
