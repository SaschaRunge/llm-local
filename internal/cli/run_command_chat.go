package cli

import _ "fmt"

func runCommandChat(ctx commandContext) error {
	ctx.cli.Mode = StateChat

	return nil
}
