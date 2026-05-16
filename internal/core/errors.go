package core

import "errors"

type ErrInvalidCommand struct {
	Context string
}

func (e ErrInvalidCommand) Error() string {
	return e.Context
}

var (
	ErrNotACommand = errors.New("not a command")
)
