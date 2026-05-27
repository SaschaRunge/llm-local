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
	return "Edit the character card for character [character_name]."
}
func (c *edit) Usage() string     { return fmt.Sprintf("%s [OPTIONAL:character_name]", c.Name()) }
func (c *edit) MinArguments() int { return 0 }
func (c *edit) MaxArguments() int { return 1 }

func (c *edit) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *edit) Execute(commandCtx core.CommandContext) (core.Result, error) {
	selection, err := selector.SelectJSONField(core.Card{}, commandCtx.Runtime.GetInput)
	if err != nil {
		return core.Result{}, err
	}

	char, err := selectCharacter(commandCtx, commandCtx.Runtime.Store().DBQueries)
	if err != nil {
		return core.Result{}, err
	}

	/*
		type user struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		}

		u := user{FirstName: "John", LastName: "Doe"}
		s := reflect.ValueOf(u)

		fmt.Println("TEEEEEST: Name:", s.FieldByName("FirstName"))*/

	newCard := core.Card{}
	json.Unmarshal(char.Card, &newCard)
	structAsReflectType := reflect.TypeOf(newCard)
	structAsValueType := reflect.ValueOf(newCard)
	index := 0
	for index = range structAsReflectType.NumField() {
		if selection == structAsReflectType.Field(index).Tag.Get("json") {
			fmt.Printf("return: %q  :   %q\n", structAsReflectType.Field(index).Name, structAsValueType.Field(index))
			break
		}
	}

	input, err := inputFromExternal(commandCtx.Runtime.Context(), fmt.Sprintf("%s", structAsValueType.Field(index)))
	if err != nil {
		return core.Result{}, err
	}
	structValue := reflect.ValueOf(&newCard).Elem()
	structFieldValue := structValue.Field(index)
	structFieldValue.SetString(strings.TrimSpace(input))

	inputAsJson, _ := json.Marshal(&newCard)
	if json.Valid(inputAsJson) {
		char.Card = json.RawMessage(inputAsJson)
		fmt.Println("YEEEEES!")
	}

	//chat, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Chat)

	err = commandCtx.Runtime.Store().DBQueries.UpdateCharacterCard(commandCtx.Runtime.Context(), database.UpdateCharacterCardParams{
		ID:   char.ID,
		Card: char.Card,
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
