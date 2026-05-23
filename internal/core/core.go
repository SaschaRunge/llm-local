package core

import (
	"context"

	"github.com/SaschaRunge/llm-local/internal/database"
)

type Result struct {
	Response  string
	NextScene Scene
}

type CommandContext struct {
	Cmd     string
	Args    []string
	Runtime Runtime
}

type Command interface {
	Name() string
	Description() string
	Usage() string
	MinAmountArguments() int
	MaxAmountArguments() int

	CanExecute(commandCtx CommandContext) bool
	Execute(commandCtx CommandContext) (Result, error)
}

type CustomParser interface {
	ParseArgs(rawArgs string) []string
}

type Runtime interface {
	Context() context.Context
	CurrentScene() Scene
	DB() *database.Queries
}

type Scene interface {
	GetName() string
	Execute(rawInput string) (Result, error)
	OnEnter() string
}
