package cli

import (
	"fmt"
	"os"

	"github.com/SaschaRunge/llm-local/internal/core"
)

func executeCommandExit(ctx commandContext) (core.Scene, error) {
	fmt.Println("=== " + Goodbye + " ===")
	os.Exit(0)

	return nil, nil
}
