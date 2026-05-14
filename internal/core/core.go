package core

import (
	"context"

	"github.com/SaschaRunge/llm-local/internal/database"
)

type Runtime interface {
	Context() context.Context
	DB() *database.Queries
	ExecuteCommand(input string) (Scene, error)
	GetInput() string
}

type Scene interface {
	Run(runtime Runtime) (Scene, error)
}
