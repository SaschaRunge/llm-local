package selector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/cli/formatters"
)

const maxAttempts = 100
const cancel = ":x"

// Selector takes a list of items []T and a conversion function toString. Selector displays the contents of []T as described by toString and allows selection through the reader function.
func Selector[T any](items []T, toString func(T) string, reader func() string) (selection T) {
	if len(items) == 0 {
		return
	}

	var response strings.Builder
	addListItem := formatters.ListBuilder(&response)

	var itemName string
	for _, item := range items {
		itemName = toString(item)
		addListItem(itemName)
	}

	fmt.Println(response.String())

	attempts := 0
	for {
		attempts++
		if attempts > maxAttempts {
			fmt.Println("No valid input provided. Smashing your keyboard doesn't solve anything. Cancel selection... .")
			return
		}

		fmt.Printf("Please type the name or an index between 1 and %d. %s to cancel.\n", len(items), cancel)
		input := strings.TrimSpace(reader())
		index, err := strconv.Atoi(input)
		if err != nil {
			for _, item := range items {
				itemName = toString(item)
				if itemName == input {
					return item
				}
			}
			if input == cancel {
				fmt.Println("Cancel selection... .")
				return
			}
			continue
		}

		index -= 1
		if index < 0 || index >= len(items) {
			fmt.Println("Please choose a valid index or name.")
			continue
		}

		return items[index]
	}
}
