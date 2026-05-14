package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
)

const (
	CmdChat        = "/chat"
	CmdDebugDelete = "/debugdelete"
)

type Cli struct {
	CommandRegistry map[string]command
	CommandAlias    CommandAlias
	Scene           core.Scene
	DBQueries       *database.Queries
	context         context.Context
	scanner         *bufio.Scanner
}

type command struct {
	name               string
	description        string
	usage              string
	callback           func(commandContext) error
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

func (c *Cli) ExecuteCommand(input string) (core.Scene, error) {
	cmd, args := parseInput(input)
	err := c.RunCommand(cmd, args)
	if err != nil {
		return nil, err
	}
	//TODO: implement logic
	return nil, nil
}

func (c *Cli) RunCommand(cmdName string, args []string) error {
	cmd, exists := c.CommandRegistry[cmdName]
	if !exists {
		return fmt.Errorf("%s is not a valid command.", cmdName)
	}

	if len(args) < cmd.minAmountArguments {
		return fmt.Errorf("Not enough arguments in %s command. Usage: %s", cmd.name, cmd.usage)
	}
	if len(args) > cmd.maxAmountArguments {
		return fmt.Errorf("To many arguments in %s command. Usage: %s", cmd.name, cmd.usage)
	}

	if cmd.requiredScene != c.Scene {
		return fmt.Errorf("Command %s not available in current context.", cmd.name)
	}

	ctx := commandContext{
		command: cmd,
		args:    args,
		cli:     c,
	}
	return cmd.callback(ctx)
}

func parseInput(input string) (cmd string, args []string) {
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
			callback:           runCommandChat,
			minAmountArguments: 1,
			maxAmountArguments: 1,
		},
		CmdDebugDelete: {
			name:               CmdDebugDelete,
			description:        "Deletes the current chat's history.",
			usage:              fmt.Sprintf("%s", CmdChat),
			callback:           runCommandDebugDelete,
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
	}
}
