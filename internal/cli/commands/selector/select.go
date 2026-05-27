package selector

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/cli/formatters"
)

const maxAttempts = 100
const cancel = ":x"

func SelectJSONField[T any](jsonStruct T, reader func() (string, error)) (selection string, e error) {
	structJSONTags := []string{}
	structAsReflectType := reflect.TypeOf(jsonStruct)
	for i := range structAsReflectType.NumField() {
		//TODO: probably use lookup
		structJSONTags = append(structJSONTags, structAsReflectType.Field(i).Tag.Get("json"))
	}

	selection, err := Select(structJSONTags, func(s string) string {
		return s
	}, reader)
	if err != nil {
		return "", err
	}

	return selection, nil
}

/*

field.Tag.Get("json")
Value)Field
Value)IsValud
(Value)Kind
(value)NumField

*/

// Select takes a list of items []T and a conversion function toString. Select displays the contents of []T as described by toString and allows selection through the reader function.
func Select[T any](items []T, toString func(T) string, reader func() (string, error)) (selection T, e error) {
	if len(items) == 0 {
		return
	}

	var response strings.Builder
	addListItem := formatters.ListBuilder(&response)

	fmt.Printf("Please type the name or an index between 1 and %d. %s to cancel.\n", len(items), cancel)

	var itemName string
	for _, item := range items {
		itemName = toString(item)
		if itemName == "" {
			itemName = "<name not found>"
		}
		addListItem(itemName)
	}

	fmt.Println(response.String())

	attempts := 0
	for {
		attempts++
		if attempts > maxAttempts {
			return selection, fmt.Errorf("No valid input provided. Smashing your keyboard doesn't solve anything. Cancel selection... .")
		}

		input, err := reader()
		if err != nil {
			return selection, err
		}

		input = strings.TrimSpace(input)
		index, err := strconv.Atoi(input)
		if err != nil {
			for _, item := range items {
				itemName = toString(item)
				if itemName == input {
					return item, nil
				}
			}
			if input == cancel {
				fmt.Println("Cancel selection... .")
				return
			}
			fmt.Printf("Specified name %q was not found. Please try again.", input)
			continue
		}

		index -= 1
		if index < 0 || index >= len(items) {
			fmt.Println("Please choose a valid index or name.")
			continue
		}

		return items[index], nil
	}
}
