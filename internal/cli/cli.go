package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

const (
	CmdChat        = "/chat"
	CmdDebugDelete = "/debugdelete"
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
	}
}

// Implements Runtime
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
			return nil, core.ErrNotACommand
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
	for {
		userInput := c.GetInput()
		currentScene, err := c.Handle(userInput)
		if err != nil {
			// TODO: implement error handling/output via cleanOutput
			fmt.Println(err.Error())
			continue
		}
		c.CurrentScene = currentScene
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
		CmdDebugDelete: {
			name:               CmdDebugDelete,
			description:        "Deletes the current chat's history.",
			usage:              fmt.Sprintf("%s", CmdChat),
			callback:           executeCommandDebugDelete,
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
	}
}
