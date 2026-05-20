package cli

import (
	"bufio"
	"context"
	_ "errors"
	"fmt"
	"os"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/cli/commands"
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

type Cli struct {
	CommandRegistry map[string]commandInfo
	CommandAlias    CommandAlias
	currentScene    core.Scene
	dbQueries       *database.Queries
	context         context.Context
	scanner         *bufio.Scanner
}

type commandInfo struct {
	name               string
	description        string
	usage              string
	command            core.Command
	requiredScene      core.Scene
	minAmountArguments int
	maxAmountArguments int
}

type CommandAlias struct {
	alias map[string][]string
}

func New(dbQueries *database.Queries) *Cli {
	return &Cli{
		CommandRegistry: getRegistry(),
		dbQueries:       dbQueries,
		scanner:         bufio.NewScanner(os.Stdin),
		context:         context.Background(),
		currentScene:    &scenes.SceneLobby{},
	}
}

// Implements core.Runtime
func (c *Cli) Context() context.Context {
	return c.context
}

func (c *Cli) CurrentScene() core.Scene {
	return c.currentScene
}

func (c *Cli) DB() *database.Queries {
	return c.dbQueries
}

func (c *Cli) GetInput() string {
	input := ""
	for {
		c.scanner.Scan()
		input = c.scanner.Text()
		if input != "" {
			break
		}
	}
	return input
}

func (c *Cli) Execute(cmdName string, args []string) (core.CommandResult, error) {
	cmdInfo, exists := c.CommandRegistry[cmdName]
	if !strings.HasPrefix(cmdName, "/") && !strings.HasPrefix(cmdName, "@") {
		cmdExecChat := &scenes.CommandExecuteChat{}

		if cmdExecChat.CanExecute(core.CommandContext{
			Cmd:     cmdName,
			Args:    args,
			Runtime: c,
		}) {
			return cmdExecChat.Execute(core.CommandContext{
				Cmd:     cmdName,
				Args:    args,
				Runtime: c,
			})
		}
	}

	if !exists {
		return core.CommandResult{}, fmt.Errorf("unknown command %q: %w", cmdName, core.ErrNotACommand)
	}

	if len(args) < cmdInfo.minAmountArguments {
		return core.CommandResult{}, core.ErrInvalidCommand{Context: fmt.Sprintf("not enough arguments in %q command. usage: %q", cmdInfo.name, cmdInfo.usage)}
	}
	if len(args) > cmdInfo.maxAmountArguments {
		return core.CommandResult{}, core.ErrInvalidCommand{Context: fmt.Sprintf("to many arguments in %q command. usage: %q", cmdInfo.name, cmdInfo.usage)}
	}

	cmdCtx := core.CommandContext{
		Args:    args,
		Runtime: c,
	}

	if !cmdInfo.command.CanExecute(cmdCtx) {
		return core.CommandResult{}, core.ErrInvalidCommand{Context: fmt.Sprintf("command %q not available in current context.", cmdInfo.name)}
	}

	return cmdInfo.command.Execute(cmdCtx)
}

func (c *Cli) Run() error {
	fmt.Println("=== " + core.Greeting + " ===")
	fmt.Printf("Entering %s:\n", c.CurrentScene().GetName())

	for {
		rawInput := c.GetInput()
		//if isCommand(rawInput) {
		cmd, args := parse(rawInput)
		result, err := c.Execute(cmd, args)
		if err != nil {
			// TODO: implement error handling/output via cleanOutput
			fmt.Println(err)
			continue
		}
		if result.NextScene != c.CurrentScene() {
			fmt.Printf("Entering %s:\n", result.NextScene.GetName())
			c.currentScene = result.NextScene
			continue
		}

		//}
		/*
			result, err := c.CurrentScene().Execute(rawInput)
			if err != nil {
				fmt.Printf("scene returned with error %q", err)
			}
		*/
		fmt.Println(result.Output)
	}
}

func isCommand(rawInput string) bool {
	return strings.HasPrefix(rawInput, "/") || strings.HasPrefix(rawInput, "@")
}

func parse(rawInput string) (cmd string, args []string) {
	parts := strings.Split(rawInput, " ")
	cmd = parts[0]
	args = []string{}
	args = append(args, parts[1:]...)

	return cmd, args
}

func translateError(err error) string {
	return err.Error()
}

func getRegistry() map[string]commandInfo {
	return map[string]commandInfo{
		CmdChat: {
			name:               CmdChat,
			description:        "Starts the chat interface as [persona]. Will load or create the chat [chat_name]. ",
			usage:              fmt.Sprintf("%s [chat_name]", CmdChat),
			command:            &commands.CommandChat{},
			minAmountArguments: 1,
			maxAmountArguments: 1,
		},
		CmdChats: {
			name:               CmdChats,
			description:        "Shows the available chats for the current user.",
			usage:              fmt.Sprintf("%s", CmdChats),
			command:            &commands.CommandChats{},
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
		CmdDebugDelete: {
			name:               CmdDebugDelete,
			description:        "Deletes the current chat's history.",
			usage:              fmt.Sprintf("%s", CmdChat),
			command:            &commands.CommandDebugDelete{},
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
		CmdExit: {
			name:               CmdExit,
			description:        "Exits the program.",
			usage:              fmt.Sprintf("%s", CmdExit),
			command:            &commands.CommandExit{},
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
		CmdLobby: {
			name:               CmdLobby,
			description:        "Enters the lobby.",
			usage:              fmt.Sprintf("%s", CmdLobby),
			command:            &commands.CommandLobby{},
			minAmountArguments: 0,
			maxAmountArguments: 0,
		},
	}
}
