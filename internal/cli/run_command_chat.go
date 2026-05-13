package cli

import (
	_ "github.com/SaschaRunge/llm-local/internal/modes"
)

func runCommandChat(ctx commandContext) error {
	ctx.cli.Mode = StateChat

	//modes.ModeChat{

	//}
	return nil
}
