package cli

import (
	"context"
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/database"
)

const (
	CmdChat        = "/chat"
	CmdDebugDelete = "/debugdelete"
)

type Cli struct {
	CommandRegistry map[string]command
	CommandAlias    CommandAlias
	Mode            State
	DBQueries       *database.Queries
}

type command struct {
	name               string
	description        string
	usage              string
	callback           func(commandContext) error
	requiredMode       State
	minAmountArguments int
	maxAmountArguments int
}

type commandContext struct {
	cli     *Cli
	command command
	args    []string
	context context.Context
}

type CommandAlias struct {
	alias map[string][]string
}

func New(dbQueries *database.Queries) *Cli {
	return &Cli{
		CommandRegistry: getRegistry(),
		Mode:            StateDefault,
		DBQueries:       dbQueries,
	}
}

func (c *Cli) RunCommand(cmdName string, args ...string) error {
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

	if cmd.requiredMode != StateAny && c.Mode != cmd.requiredMode {
		return fmt.Errorf("Unable to run %s command in current context.", cmd.name)
	}

	ctx := commandContext{
		cli:     c,
		command: cmd,
		args:    args,
		context: context.Background(),
	}
	return cmd.callback(ctx)
}

func getRegistry() map[string]command {
	return map[string]command{
		CmdChat: {
			name:               CmdChat,
			description:        "Starts the chat interface. If a [chat_name] argument is provided, will load or create the chat.",
			usage:              fmt.Sprintf("%s [chat_name]", CmdChat),
			callback:           runCommandChat,
			requiredMode:       StateAny,
			minAmountArguments: 1,
			maxAmountArguments: 1,
		},
		CmdDebugDelete: {
			name:               CmdDebugDelete,
			description:        "Deletes the current chat's history.",
			usage:              fmt.Sprintf("%s", CmdChat),
			callback:           runCommandDebugDelete,
			requiredMode:       StateChat,
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
	}
}
