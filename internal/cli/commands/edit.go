package commands

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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

	fmt.Printf("Please enter a %s:\n", selectionName)
	input, err := commandCtx.Runtime.GetInput(fmt.Sprintf("%s", cardFieldValue))
	if err != nil {
		return core.Result{}, err
	}

	switch cardFieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		inputAsInt, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return core.Result{}, err
		}
		cardFieldValue.SetInt(int64(inputAsInt))
	case reflect.String:
		cardFieldValue.SetString(strings.TrimSpace(input))
	default:
		return core.Result{}, fmt.Errorf("unsupported field type in json data")

	}

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
