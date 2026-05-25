package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/SaschaRunge/llm-local/internal/cli/commands"
	"github.com/SaschaRunge/llm-local/internal/cli/parser"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/scenes"
)

type Cli struct {
	CommandRegistry map[string]core.Command
	CommandAlias    CommandAlias
	currentScene    core.Scene
	context         context.Context
	scanner         *bufio.Scanner
	store           *core.Store
}

type CommandAlias struct {
	alias map[string][]string
}

func New(store *core.Store) *Cli {
	return &Cli{
		CommandRegistry: getRegistry(),
		context:         context.Background(),
		scanner:         bufio.NewScanner(os.Stdin),
		store:           store,
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

func (c *Cli) Store() *core.Store {
	return c.store
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

func (c *Cli) Execute(cmd core.Command, rawArgs string) (core.Result, error) {
	var args []string
	switch customParser := cmd.(type) {
	case core.CustomParser:
		args = customParser.ParseArgs(rawArgs)
	default:
		args = parser.ParseArgs(rawArgs)
	}

	if len(args) < cmd.MinArguments() {
		return core.Result{}, core.ErrInvalidCommand{Context: fmt.Sprintf("not enough arguments in %q command. usage: %q", cmd.Name(), cmd.Usage())}
	}
	if len(args) > cmd.MaxArguments() {
		return core.Result{}, core.ErrInvalidCommand{Context: fmt.Sprintf("too many arguments in %q command. usage: %q", cmd.Name(), cmd.Usage())}
	}

	cmdCtx := core.CommandContext{
		Args:    args,
		Runtime: c,
	}

	if !cmd.CanExecute(cmdCtx) {
		return core.Result{}, core.ErrInvalidCommand{Context: fmt.Sprintf("command %q not available in scene %q.", cmd.Name(), c.CurrentScene().GetName())}
	}

	return cmd.Execute(cmdCtx)
}

func (c *Cli) Run() error {
	fmt.Println("=== " + core.Greeting + " ===")
	fmt.Printf("Entering %s:\n\n", c.CurrentScene().GetName())

	for {
		var result core.Result
		var err error

		rawInput := c.GetInput()
		preprocessor, inputIsCommand := parser.SelectPreprocessorByPrefix(rawInput)
		if inputIsCommand {
			cmd, rawArgs := preprocessor(rawInput)
			result, err = c.handleCommand(cmd, rawArgs)
		} else {
			result, err = c.handleRawInput(rawInput)
		}
		if err != nil {
			fmt.Printf("unable to process input, error: %s", err)
		}

		if result.NextScene != nil && result.NextScene != c.CurrentScene() {
			fmt.Printf("Entering %s:\n\n", result.NextScene.GetName())
			c.currentScene = result.NextScene
		} else {
			fmt.Println(result.Response)
		}
	}
}

/* //TODO: do something on enter/exit
func (c *Cli) sceneChange(nextScene core.Scene) {
	nextScene.OnEnter()
	c.currentScene.OnExit()
	c.currentScene = nextScene
}
*/

func (c *Cli) handleRawInput(rawInput string) (core.Result, error) {
	var err error
	result := core.Result{}

	if scene, ok := c.CurrentScene().(core.AllowsRawInput); ok {
		result, err = scene.HandleRawInput(rawInput)
		if err != nil {
			err = fmt.Errorf("handling input failed with error: %w", err)
		}
	} else {
		err = fmt.Errorf("please enter a valid command")
	}

	return result, err
}

func (c *Cli) handleCommand(cmd, rawArgs string) (core.Result, error) {
	var err error
	result := core.Result{}

	if _, exists := c.CommandRegistry[cmd]; !exists {
		err = fmt.Errorf("unknown command %q: %w", cmd, core.ErrNotACommand)
	} else {
		result, err = c.Execute(c.CommandRegistry[cmd], rawArgs)
	}

	return result, err
}

func getRegistry() map[string]core.Command {
	registry := make(map[string]core.Command)
	for _, cmd := range commands.All() {
		registry[cmd.Name()] = cmd
	}
	return registry
}
