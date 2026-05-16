package core

import (
	"context"

	"github.com/SaschaRunge/llm-local/internal/database"
)

type Runtime interface {
	Context() context.Context
	DB() *database.Queries
	Handle(input string) (Scene, error)
}

type Scene interface {
	Execute(cmd string) (Scene, error)
}
