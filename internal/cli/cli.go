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

func GetPrefixFrom(rawInput string) string {
	validPrefixes := []string{
		"/",
		"@",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(rawInput, prefix) {
			return prefix
		}
	}

	return ""
}

func (c *Cli) Execute(cmd core.Command, rawArgs string) (core.Result, error) {
	var args []string
	switch parser := cmd.(type) {
	case core.CustomParser:
		args = parser.ParseArgs(rawArgs)
		//fmt.Printf("%v", args)
	default:
		args = parse((rawArgs))
	}

	if len(args) < cmd.MinAmountArguments() {
		return core.Result{}, core.ErrInvalidCommand{Context: fmt.Sprintf("not enough arguments in %q command. usage: %q", cmd.Name(), cmd.Usage())}
	}
	if len(args) > cmd.MaxAmountArguments() {
		return core.Result{}, core.ErrInvalidCommand{Context: fmt.Sprintf("to many arguments in %q command. usage: %q", cmd.Name(), cmd.Usage())}
	}

	cmdCtx := core.CommandContext{
		Args:    args,
		Runtime: c,
	}

	if !cmd.CanExecute(cmdCtx) {
		return core.Result{}, core.ErrInvalidCommand{Context: fmt.Sprintf("command %q not available in current context.", cmd.Name())}
	}

	return cmd.Execute(cmdCtx)
}

func (c *Cli) Run() error {
	fmt.Println("=== " + core.Greeting + " ===")
	fmt.Printf("Entering %s:\n", c.CurrentScene().GetName())

	//TODO: clean up flow
	for {
		rawInput := c.GetInput()
		//if isCommand(rawInput) {
		var result core.Result
		var err error
		if GetPrefixFrom(rawInput) == "" {
			//TODO: split into optional interface
			result, err = c.currentScene.Execute(rawInput)
			if err != nil {
				// TODO: implement error handling/output via cleanOutput
				fmt.Println(err)
				continue
			}
		} else {
			cmd, rawArgs := extractCommandString(rawInput)
			if _, exists := c.CommandRegistry[cmd]; !exists {
				fmt.Printf("unknown command %q: %q", cmd, core.ErrNotACommand)
			} else {
				result, err = c.Execute(c.CommandRegistry[cmd], rawArgs)
				if err != nil {
					// TODO: implement error handling/output via cleanOutput
					fmt.Println(err)
					continue
				}
			}
		}

		if result.NextScene != c.CurrentScene() {
			fmt.Printf("Entering %s:\n", result.NextScene.GetName())
			c.currentScene = result.NextScene
			continue
		}

		fmt.Println(result.Response)
	}
}

func extractCommandString(rawInput string) (cmd string, args string) {
	parts := strings.SplitN(rawInput, " ", 2)

	cmd = parts[0]
	if len(parts) >= 2 {
		args = parts[1]
	}

	return cmd, args
}

func parse(rawArgs string) (args []string) {
	if rawArgs == "" {
		return args
	}
	return strings.Split(rawArgs, " ")
}

func getRegistry() map[string]core.Command {
	registry := make(map[string]core.Command)
	for _, cmd := range (&commands.Library{}).GetAll() {
		registry[cmd.Name()] = cmd
	}
	return registry
}
