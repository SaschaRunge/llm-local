package cli

import "fmt"

type Cli struct {
	Commands     map[string]command
	CommandAlias CommandAlias
	Mode         Mode
}

type command struct {
	name               string
	description        string
	usage              string
	callback           func(commandContext) error
	minAmountArguments int
	maxAmountArguments int
}

type commandContext struct {
	cli     *Cli
	command command
	args    []string
}

type CommandAlias struct {
	alias map[string][]string
}

func New() *Cli {
	return &Cli{
		Commands: getCommands(),
	}
}

func (c *Cli) RunCommand(cmdName string, args ...string) error {
	cmd, exists := c.Commands[cmdName]
	if !exists {
		return fmt.Errorf("%s is not a valid command.", cmdName)
	}

	ctx := commandContext{
		cli:     c,
		command: cmd,
		args:    args,
	}
	return cmd.callback(ctx)
}

func getCommands() map[string]command {
	return map[string]command{
		"/chat": {
			name:               "chat",
			description:        "Starts the chat interface. If a [chat_name] argument is provided, will load or create the chat.",
			usage:              "/chat [chat_name]",
			callback:           commandChat,
			minAmountArguments: 1,
			maxAmountArguments: 2,
		},
	}
}
