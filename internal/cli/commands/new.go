package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/SaschaRunge/llm-local/internal/scenes"

	"github.com/google/uuid"
)

var DefaultUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type new struct{}

func (c *new) Name() string { return "/new" }
func (c *new) Description() string {
	return "Create a new chat or character"
}
func (c *new) Usage() string     { return fmt.Sprintf("%s [option] [name]", c.Name()) }
func (c *new) MinArguments() int { return 2 }
func (c *new) MaxArguments() int { return 2 }

func (c *new) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (c *new) Execute(commandCtx core.CommandContext) (core.Result, error) {
	var options = map[string]option{
		"chat":      newChat,
		"character": newCharacter,
	}

	optionArg := commandCtx.Args[0]
	executeOption, ok := options[optionArg]
	if !ok {
		return core.Result{}, fmt.Errorf("%q is not a valid option for command %q. usage: %s", optionArg, c.Name(), c.Usage())
	}

	return executeOption(commandCtx)
}

func (c *new) ParseArgs(rawArgs string) []string {
	parts := strings.SplitN(rawArgs, " ", 2)
	return parts
}

func newChat(commandCtx core.CommandContext) (core.Result, error) {
	var newChat database.Chat
	var selectedCharacter character
	var err error
	err = commandCtx.Runtime.Store().ExecTx(commandCtx.Runtime.Context(), func(qtx *database.Queries) error {
		selectedCharacter, err = selectCharacter(commandCtx, qtx)
		if err != nil {
			return err
		}

		newChat, err = qtx.AddChat(commandCtx.Runtime.Context(), database.AddChatParams{
			Name:            strings.TrimSpace(commandCtx.Args[1]),
			UserCharacterID: DefaultUserID,
		})
		if err != nil {
			return err
		}

		err = qtx.SubscribeToChat(commandCtx.Runtime.Context(), database.SubscribeToChatParams{
			ChatID:      newChat.ID,
			CharacterID: selectedCharacter.ID,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return core.Result{}, err
	}

	var response strings.Builder
	fmt.Fprintf(&response, "Chat %q successfully created.\n", newChat.Name)
	fmt.Fprintf(&response, "Character %q added.", selectedCharacter.Name)

	return core.Result{Response: response.String()}, nil
}

func newCharacter(commandCtx core.CommandContext) (core.Result, error) {
	var err error
	var result core.Result
	err = commandCtx.Runtime.Store().ExecTx(commandCtx.Runtime.Context(), func(qtx *database.Queries) error {
		newCharacter, err := qtx.AddCharacter(commandCtx.Runtime.Context(), database.AddCharacterParams{
			Name: strings.TrimSpace(commandCtx.Args[1]),
			Card: json.RawMessage("{}"),
			//TODO: allow for creation of user character. maybe own option user
			IsUser: false,
		})
		if err != nil {
			return err
		}

		switch currentScene := commandCtx.Runtime.CurrentScene().(type) {
		case *scenes.Chat:
			err := qtx.SubscribeToChat(commandCtx.Runtime.Context(), database.SubscribeToChatParams{
				ChatID:      currentScene.ID,
				CharacterID: newCharacter.ID,
			})
			if err != nil {
				return err
			}

			result = core.Result{Response: fmt.Sprintf("Character %q successfully created and subscribed to Chat %q.", newCharacter.Name, currentScene.GetName())}
			return nil
		default:
			result = core.Result{Response: fmt.Sprintf("Character %q successfully created.", newCharacter.Name)}
			return nil
		}
	})
	if err != nil {
		return core.Result{}, err
	}

	return result, nil
}
