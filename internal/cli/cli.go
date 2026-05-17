package cli

import (
	"bufio"
	"context"
	_ "errors"
	"fmt"
	"os"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

const (
	CmdChat        = "/chat"
	CmdChats       = "/chats"
	CmdDebugDelete = "/debugdelete"
	CmdExit        = "/exit"
	CmdLobby       = "/lobby"
)

const (
	Greeting = "Welcome to llm-local."
	Goodbye  = "Closing... . Goodbye."
)

type Cli struct {
	CommandRegistry map[string]command
	CommandAlias    CommandAlias
	CurrentScene    core.Scene
	DBQueries       *database.Queries
	context         context.Context
	scanner         *bufio.Scanner
}

type command struct {
	name               string
	description        string
	usage              string
	callback           func(commandContext) (core.Scene, error)
	requiredScene      core.Scene
	minAmountArguments int
	maxAmountArguments int
}

type commandContext struct {
	command command
	args    []string
	cli     *Cli
}

type CommandAlias struct {
	alias map[string][]string
}

func New(dbQueries *database.Queries) *Cli {
	return &Cli{
		CommandRegistry: getRegistry(),
		DBQueries:       dbQueries,
		scanner:         bufio.NewScanner(os.Stdin),
		context:         context.Background(),
		CurrentScene:    &scenes.SceneLobby{},
	}
}

// Implements core.Runtime
func (c *Cli) Context() context.Context {
	return c.context
}

func (c *Cli) DB() *database.Queries {
	return c.DBQueries
}

func (c *Cli) GetInput() string {
	c.scanner.Scan()
	return c.scanner.Text()
}

// TODO: probably should split
func (c *Cli) Handle(input string) (core.Scene, error) {
	if !strings.HasPrefix(input, "/") {
		switch c.CurrentScene.(type) {
		case *scenes.SceneChat:
			return c.CurrentScene, nil
		default:
			return nil, fmt.Errorf("unknown command %q: %w", input, core.ErrNotACommand)
		}
	}

	cmd, args := parse(input)

	return c.Execute(cmd, args)
}

func (c *Cli) Execute(cmdName string, args []string) (core.Scene, error) {
	cmd, exists := c.CommandRegistry[cmdName]
	if !exists {
		return nil, fmt.Errorf("unknown command %q: %w", cmdName, core.ErrNotACommand)
	}

	if len(args) < cmd.minAmountArguments {
		return nil, core.ErrInvalidCommand{Context: fmt.Sprintf("not enough arguments in %q command. usage: %q", cmd.name, cmd.usage)}
	}
	if len(args) > cmd.maxAmountArguments {
		return nil, core.ErrInvalidCommand{Context: fmt.Sprintf("to many arguments in %q command. usage: %q", cmd.name, cmd.usage)}
	}

	if cmd.requiredScene != nil && cmd.requiredScene != c.CurrentScene {
		return nil, core.ErrInvalidCommand{Context: fmt.Sprintf("command %q not available in current context.", cmd.name)}
	}

	ctx := commandContext{
		command: cmd,
		args:    args,
		cli:     c,
	}
	return cmd.callback(ctx)
}

func (c *Cli) Run() error {
	fmt.Println("=== " + Greeting + " ===")
	fmt.Printf("Entering %s:\n", c.CurrentScene.GetName())

	for {
		userInput := c.GetInput()
		nextScene, err := c.Handle(userInput)
		if err != nil {
			// TODO: implement error handling/output via cleanOutput
			fmt.Println(err)
			continue
		}
		if nextScene != c.CurrentScene {
			fmt.Printf("Entering %s:\n", nextScene.GetName())
		}
		c.CurrentScene = nextScene
	}
}

// TODO: add logic (do i even need?)
func TranslateError(err error) string {
	return err.Error()
}

func parse(input string) (cmd string, args []string) {
	parts := strings.Split(input, " ")
	cmd = parts[0]
	args = []string{}
	args = append(args, parts[1:]...)

	return cmd, args
}

func getRegistry() map[string]command {
	return map[string]command{
		CmdChat: {
			name:               CmdChat,
			description:        "Starts the chat interface. If a [chat_name] argument is provided, will load or create the chat.",
			usage:              fmt.Sprintf("%s [chat_name]", CmdChat),
			callback:           executeCommandChat,
			minAmountArguments: 1,
			maxAmountArguments: 1,
		},
		CmdChats: {
			name:               CmdChats,
			description:        "Shows the available chats for the current user.",
			usage:              fmt.Sprintf("%s", CmdChats),
			callback:           executeCommandChats,
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
		CmdDebugDelete: {
			name:               CmdDebugDelete,
			description:        "Deletes the current chat's history.",
			usage:              fmt.Sprintf("%s", CmdChat),
			callback:           executeCommandDebugDelete,
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
		CmdExit: {
			name:               CmdExit,
			description:        "Exits the program.",
			usage:              fmt.Sprintf("%s", CmdExit),
			callback:           executeCommandExit,
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
		CmdLobby: {
			name:               CmdLobby,
			description:        "Enters the lobby.",
			usage:              fmt.Sprintf("%s", CmdLobby),
			callback:           executeCommandLobby,
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
	}
}
