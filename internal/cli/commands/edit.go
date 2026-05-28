package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/cli/commands/selector"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
)

type edit struct{}

func (c *edit) Name() string { return "/edit" }
func (c *edit) Description() string {
	return "Edit the card for the selected chat or character."
}
func (c *edit) Usage() string     { return fmt.Sprintf("%s", c.Name()) }
func (c *edit) MinArguments() int { return 0 }
func (c *edit) MaxArguments() int { return 0 }

func (c *edit) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *edit) Execute(commandCtx core.CommandContext) (core.Result, error) {
	selectedCharacter, err := selectCharacter(commandCtx, commandCtx.Runtime.Store().DBQueries)
	if err != nil {
		return core.Result{}, err
	}

	selectionName, selectionPosition, err := selector.SelectJSONField(core.Card{}, commandCtx.Runtime.GetInput)
	if err != nil {
		return core.Result{}, err
	}

	newCard := core.Card{}
	json.Unmarshal(selectedCharacter.Card, &newCard)

	cardValue := reflect.ValueOf(&newCard).Elem()
	cardFieldValue := cardValue.Field(selectionPosition)

	if !cardFieldValue.IsValid() {
		return core.Result{}, fmt.Errorf("unexpected error: position %d doesn't exist in card struct", selectionPosition)
	}

	if !cardFieldValue.CanSet() {
		return core.Result{}, fmt.Errorf("cannot set %s field value", selectionName)
	}

	input, err := inputFromExternal(commandCtx.Runtime.Context(), fmt.Sprintf("%s", cardFieldValue))
	if err != nil {
		return core.Result{}, err
	}

	cardFieldValue.SetString(strings.TrimSpace(input))

	inputAsJSON, err := json.Marshal(&newCard)
	if err != nil {
		return core.Result{}, err
	}
	if !json.Valid(inputAsJSON) {
		return core.Result{}, fmt.Errorf("%s is not valid json", inputAsJSON)
	}

	selectedCharacter.Card = json.RawMessage(inputAsJSON)

	err = commandCtx.Runtime.Store().DBQueries.UpdateCharacterCard(commandCtx.Runtime.Context(), database.UpdateCharacterCardParams{
		ID:   selectedCharacter.ID,
		Card: selectedCharacter.Card,
	})
	if err != nil {
		return core.Result{}, err
	}

	return core.Result{}, nil
}

func inputFromExternal(ctx context.Context, data string) (string, error) {
	const workingDir = "./"

	workingDirAbs, _ := filepath.Abs(workingDir)

	tmp, err := os.CreateTemp(workingDirAbs, "tmp_input_")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	os.WriteFile(tmp.Name(), []byte(data), 2)
	tmp.Close()

	//likely shouldn't be with context so the user doesn't lose his input on crash
	cmd := exec.CommandContext(ctx, "nvim", tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	input, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", err
	}

	return string(input), nil
}
