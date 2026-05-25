package core

import (
	"context"
	"database/sql"

	"github.com/SaschaRunge/llm-local/internal/database"
)

type CommandContext struct {
	Cmd     string
	Args    []string
	Runtime Runtime
}

type Result struct {
	Response  string
	NextScene Scene
}

type Store struct {
	DB        *sql.DB
	DBQueries *database.Queries
}

type AllowsRawInput interface {
	HandleRawInput(rawInput string) (Result, error)
}

type Command interface {
	Name() string
	Description() string
	Usage() string
	MinArguments() int
	MaxArguments() int

	CanExecute(commandCtx CommandContext) bool
	Execute(commandCtx CommandContext) (Result, error)
}

type CustomParser interface {
	ParseArgs(rawArgs string) []string
}

type Runtime interface {
	Context() context.Context
	CurrentScene() Scene
	GetInput() string
	Store() *Store
}

type Scene interface {
	GetName() string
	OnEnter() string
}
