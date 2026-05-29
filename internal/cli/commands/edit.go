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
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type edit struct{}

func (c *edit) Name() string { return "/edit" }
func (c *edit) Description() string {
	return "Edit the card for the selected chat or character."
}
func (c *edit) Usage() string     { return fmt.Sprintf("%s [option] [name]", c.Name()) }
func (c *edit) MinArguments() int { return 1 }
func (c *edit) MaxArguments() int { return 1 }

func (c *edit) CanExecute(commandCtx core.CommandContext) bool {
	//TODO: temporary, because otherwise i'd have to refresh the chat when calling from chat
	_, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Lobby)
	return isSceneChat
}

func (c *edit) Execute(commandCtx core.CommandContext) (core.Result, error) {
	var options = map[string]option{
		"chat":      editChat,
		"character": editCharacter,
	}

	optionArg := commandCtx.Args[0]
	executeOption, ok := options[optionArg]
	if !ok {
		return core.Result{}, fmt.Errorf("%q is not a valid option for command %q. usage: %s", optionArg, c.Name(), c.Usage())
	}

	return executeOption(commandCtx)

}

func editChat(commandCtx core.CommandContext) (core.Result, error) {
	selectedChat, err := selectChat(commandCtx, commandCtx.Runtime.Store().DBQueries)
	if err != nil {
		return core.Result{}, err
	}

	var changedTag string
	selectedChat.Card, changedTag, err = editCard(commandCtx, selectedChat.Card)
	if err != nil {
		return core.Result{}, err
	}

	insertedChatID, err := commandCtx.Runtime.Store().DBQueries.UpdateChatCard(commandCtx.Runtime.Context(), database.UpdateChatCardParams{
		ID:   selectedChat.ID,
		Card: selectedChat.Card,
	})
	if err != nil {
		return core.Result{}, err
	}
	if insertedChatID != selectedChat.ID {
		return core.Result{}, fmt.Errorf("unexpected error: can't find correspondig id for inserting chat after udpate")
	}

	return core.Result{Response: fmt.Sprintf("\n%s successfully saved!", changedTag)}, nil
}

func editCharacter(commandCtx core.CommandContext) (core.Result, error) {
	selectedCharacter, err := selectCharacter(commandCtx, commandCtx.Runtime.Store().DBQueries)
	if err != nil {
		return core.Result{}, err
	}

	var changedTag string
	selectedCharacter.Card, changedTag, err = editCard(commandCtx, selectedCharacter.Card)
	if err != nil {
		return core.Result{}, err
	}

	insertedCharacterID, err := commandCtx.Runtime.Store().DBQueries.UpdateCharacterCard(commandCtx.Runtime.Context(), database.UpdateCharacterCardParams{
		ID:   selectedCharacter.ID,
		Card: selectedCharacter.Card,
	})
	if err != nil {
		return core.Result{}, err
	}
	if insertedCharacterID != selectedCharacter.ID {
		return core.Result{}, fmt.Errorf("unexpected error: can't find correspondig id for inserting character after udpate")
	}

	return core.Result{Response: fmt.Sprintf("\n%s successfully saved!", changedTag)}, nil
}

func editCard(commandCtx core.CommandContext, card json.RawMessage) (changedCard json.RawMessage, jsonTag string, e error) {
	selectionName, selectionPosition, err := selector.SelectJSONField(core.Card{}, commandCtx.Runtime.GetInput)
	if err != nil {
		return nil, "", err
	}

	newCard := core.Card{}
	json.Unmarshal(card, &newCard)

	cardValue := reflect.ValueOf(&newCard).Elem()
	cardFieldValue := cardValue.Field(selectionPosition)

	if !cardFieldValue.IsValid() {
		return nil, "", fmt.Errorf("unexpected error: position %d doesn't exist in card struct", selectionPosition)
	}

	if !cardFieldValue.CanSet() {
		return nil, "", fmt.Errorf("cannot set %s field value", selectionName)
	}

	fmt.Printf("Please enter a %s:\n", selectionName)
	input, err := commandCtx.Runtime.GetInput(fmt.Sprintf("%s", cardFieldValue))
	if err != nil {
		return nil, "", err
	}

	switch cardFieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		inputAsInt, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return nil, "", err
		}
		//TODO: conversion of int is not safe, gets cast down if field doesnt have enough byte size
		cardFieldValue.SetInt(int64(inputAsInt))
	case reflect.String:
		cardFieldValue.SetString(strings.TrimSpace(input))
	default:
		return nil, "", fmt.Errorf("unsupported field type in json data")
	}

	inputAsJSON, err := json.Marshal(&newCard)
	if err != nil {
		return nil, "", err
	}
	if !json.Valid(inputAsJSON) {
		return nil, "", fmt.Errorf("%s is not valid json", inputAsJSON)
	}

	return json.RawMessage(inputAsJSON), selectionName, nil
}
