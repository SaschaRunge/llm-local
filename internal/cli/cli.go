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

type Cli struct {
	CommandRegistry map[string]core.Command
	CommandAlias    CommandAlias
	currentScene    core.Scene
	dbQueries       *database.Queries
	context         context.Context
	scanner         *bufio.Scanner
}

type commandInfo struct {
	name        string
	description string
	usage       string
	//command            core.Command
	//requiredScene      core.Scene
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
		currentScene:    &scenes.Lobby{},
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
	cmd, exists := c.CommandRegistry[cmdName]

	if !strings.HasPrefix(cmdName, "/") && !strings.HasPrefix(cmdName, "@") {
		//TODO: implement fallback via registry
		cmdExecChat := &scenes.CommandExecuteChat{}

		if cmdExecChat.CanExecute(core.CommandContext{
			//Cmd:     cmdName + strings.Join(args, " "),
			//Args:    args,
			Runtime: c,
		}) {
			return cmdExecChat.Execute(core.CommandContext{
				Cmd: cmdName + " " + strings.Join(args, " "),
				//Args:    args,
				Runtime: c,
			})
		}
	}

	if !exists {
		return core.CommandResult{}, fmt.Errorf("unknown command %q: %w", cmdName, core.ErrNotACommand)
	}

	if len(args) < cmd.MinAmountArguments() {
		return core.CommandResult{}, core.ErrInvalidCommand{Context: fmt.Sprintf("not enough arguments in %q command. usage: %q", cmd.Name(), cmd.Usage())}
	}
	if len(args) > cmd.MaxAmountArguments() {
		return core.CommandResult{}, core.ErrInvalidCommand{Context: fmt.Sprintf("to many arguments in %q command. usage: %q", cmd.Name(), cmd.Usage())}
	}

	cmdCtx := core.CommandContext{
		Args:    args,
		Runtime: c,
	}

	if !cmd.CanExecute(cmdCtx) {
		return core.CommandResult{}, core.ErrInvalidCommand{Context: fmt.Sprintf("command %q not available in current context.", cmd.Name())}
	}

	return cmd.Execute(cmdCtx)
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

func parse(rawInput string) (cmd string, args []string) {
	parts := strings.Split(rawInput, " ")
	cmd = parts[0]
	args = []string{}
	args = append(args, parts[1:]...)

	return cmd, args
}

func getRegistry() map[string]core.Command {
	registry := make(map[string]core.Command)
	for _, cmd := range (&commands.Library{}).GetAll() {
		registry[cmd.Name()] = cmd
	}
	return registry
}
