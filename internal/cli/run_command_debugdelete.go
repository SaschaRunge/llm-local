package cli

import "fmt"

func runCommandDebugDelete(ctx commandContext) error {
	err := ctx.cli.DBQueries.DeleteMessages(ctx.context)
	if err != nil {
		return err
	}

	fmt.Println("Debug Delete Messages.")
	return nil
}
