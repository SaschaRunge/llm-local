package core

import (
	"context"

	"github.com/SaschaRunge/llm-local/internal/database"
)

type SceneResult struct {
	Response  string
	NextScene Scene
}

type Runtime interface {
	Context() context.Context
	DB() *database.Queries
	Handle(input string) (Scene, error)
}

type Scene interface {
	GetName() string
	Execute(cmd string) (SceneResult, error)
}
