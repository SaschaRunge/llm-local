package cli

import "fmt"

func commandChat(ctx commandContext) error {
	cmd := ctx.command
	args := ctx.args
	if len(args) < ctx.command.minAmountArguments {
		return fmt.Errorf("Not enough arguments in %s command. Usage: %s", cmd.name, cmd.usage)
	}
	if len(args) > ctx.command.maxAmountArguments {
		return fmt.Errorf("To many arguments in %s command. Usage: %s", cmd.name, cmd.usage)
	}

	ctx.cli.Mode = ModeChat

	return nil
}
