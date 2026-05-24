package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"

	"github.com/google/uuid"
)

var DefaultUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type new struct{}

type option func(core.CommandContext) (core.Result, error)

var options = map[string]option{
	"chat":      optionChat,
	"character": optionCharacter,
}

func (n *new) Name() string { return "/new" }
func (n *new) Description() string {
	return "Create a new chat or character"
}
func (n *new) Usage() string     { return fmt.Sprintf("%s [option] [name]", n.Name()) }
func (n *new) MinArguments() int { return 2 }
func (n *new) MaxArguments() int { return 2 }

func (n *new) CanExecute(commandCtx core.CommandContext) bool {
	return true
}

func (n *new) Execute(commandCtx core.CommandContext) (core.Result, error) {
	optionArg := commandCtx.Args[0]
	executeOption, ok := options[optionArg]
	if !ok {
		return core.Result{}, fmt.Errorf("%q is not a valid option for command %q. usage: %s", optionArg, n.Name(), n.Usage())
	}

	return executeOption(commandCtx)
}

func optionChat(commandCtx core.CommandContext) (core.Result, error) {
	newChat, err := commandCtx.Runtime.DB().AddChat(commandCtx.Runtime.Context(), database.AddChatParams{
		Name:            commandCtx.Args[1],
		UserCharacterID: DefaultUserID,
	})
	if err != nil {
		return core.Result{}, err
	}

	return core.Result{Response: fmt.Sprintf("Chat %q successfully created.", newChat.Name)}, nil
}

func optionCharacter(commandCtx core.CommandContext) (core.Result, error) {
	/*

		switch scene := commandCtx.Runtime.CurrentScene().(type){
		case *scenes.Chat:
			commandCtx.Runtime.DB().SubscribeToChat(commandCtx.Runtime.Context(), database.SubscribeToChatParams{

			})
			scene.ID
		}*/
	return core.Result{}, nil
}
