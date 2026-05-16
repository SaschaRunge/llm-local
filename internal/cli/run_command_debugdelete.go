package cli

import (
	"fmt"
	"github.com/SaschaRunge/llm-local/internal/core"
)

func runCommandDebugDelete(ctx commandContext) (core.Scene, error) {
	err := ctx.cli.DBQueries.DeleteMessages(ctx.cli.context)
	if err != nil {
		return nil, err
	}

	fmt.Println("Debug Delete Messages.")
	return nil, nil
}
