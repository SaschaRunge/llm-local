package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type add struct{}

func (c *add) Name() string { return "/add" }
func (c *add) Description() string {
	return "Adds a character to the current chat."
}
func (c *add) Usage() string     { return fmt.Sprintf("%s", c.Name()) }
func (c *add) MinArguments() int { return 0 }
func (c *add) MaxArguments() int { return 0 }

func (c *add) CanExecute(commandCtx core.CommandContext) bool {
	_, isSceneChat := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	return isSceneChat
}

func (c *add) Execute(commandCtx core.CommandContext) (core.Result, error) {
	chat, ok := commandCtx.Runtime.CurrentScene().(*scenes.Chat)
	if !ok {
		return core.Result{}, fmt.Errorf("unexpected error in command /add: type assertion failed")
	}

	var selectedCharacter character
	var err error
	err = commandCtx.Runtime.Store().ExecTx(commandCtx.Runtime.Context(), func(qtx *database.Queries) error {
		selectedCharacter, err = selectCharacter(commandCtx, qtx)
		if err != nil {
			return err
		}

		err = qtx.SubscribeToChat(commandCtx.Runtime.Context(), database.SubscribeToChatParams{
			ChatID:      chat.ID,
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

	return core.Result{Response: fmt.Sprintf("Character %q successfully subscribed to chat %q.", selectedCharacter.Name, chat.GetName())}, nil
}
