package commands

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/SaschaRunge/llm-local/internal/scenes"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

var DefaultUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type new struct{}

type option func(core.CommandContext) (core.Result, error)

var options = map[string]option{
	"chat":      optionChat,
	"character": optionCharacter,
}

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

func optionChat(commandCtx core.CommandContext) (core.Result, error) {
	newChat, err := commandCtx.Runtime.DB().AddChat(commandCtx.Runtime.Context(), database.AddChatParams{
		Name:            strings.TrimSpace(commandCtx.Args[1]),
		UserCharacterID: DefaultUserID,
	})
	if err != nil {
		return core.Result{}, err
	}

	return core.Result{Response: fmt.Sprintf("Chat %q successfully created.", newChat.Name)}, nil
}

func optionCharacter(commandCtx core.CommandContext) (core.Result, error) {
	newCharacter, err := commandCtx.Runtime.DB().AddCharacter(commandCtx.Runtime.Context(), database.AddCharacterParams{
		Name: strings.TrimSpace(commandCtx.Args[1]),
		Card: pqtype.NullRawMessage{},
		//TODO: allow for creation of user character. maybe own option user
		IsUser: sql.NullBool{Bool: false, Valid: true},
	})
	if err != nil {
		return core.Result{}, err
	}

	switch currentScene := commandCtx.Runtime.CurrentScene().(type) {
	case *scenes.Chat:
		err := commandCtx.Runtime.DB().SubscribeToChat(commandCtx.Runtime.Context(), database.SubscribeToChatParams{
			ChatID:      currentScene.ID,
			CharacterID: newCharacter.ID,
		})
		if err != nil {
			return core.Result{}, err
		}

		return core.Result{Response: fmt.Sprintf("Character %q successfully created and subscribed to Chat %q.", newCharacter.Name, currentScene.GetName())}, nil
	default:
		return core.Result{Response: fmt.Sprintf("Character %q successfully created.", newCharacter.Name)}, nil
	}
}
