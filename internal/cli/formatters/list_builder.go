package formatters

import (
	"fmt"
	"strings"
)

func ListBuilder(stringBuilder *strings.Builder) func(string) {
	i := 0

	return func(item string) {
		fmt.Fprintf(stringBuilder, "  %d. %s\n", i+1, item)
		i++
	}
}
