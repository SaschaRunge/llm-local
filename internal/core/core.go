package core

import (
	"context"
	"database/sql"

	"github.com/SaschaRunge/llm-local/internal/database"
)

// TODO: currently duplicate of chat.card
type Card struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Personality string `json:"personality"`
	FirstMsg    string `json:"first_msg"`
	MsgExample  string `json:"msg_example"`
	NotAString  int    `json:"age"`
}

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

func (s *Store) ExecTx(ctx context.Context, fn func(*database.Queries) error) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.DBQueries.WithTx(tx)
	err = fn(qtx)
	if err != nil {
		return err
	}

	return tx.Commit()
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
	GetInput(string) (string, error)
	Store() *Store
}

type Scene interface {
	GetName() string
	OnEnter() string
}
