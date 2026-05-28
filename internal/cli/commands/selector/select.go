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

func SelectJSONField[T any](jsonStruct T, reader func() (string, error)) (selectionName string, selectionIndex int, e error) {
	jsonTags := []string{}
	jsonTagPositions := map[string]int{}

	jsonStructReflectType := reflect.TypeOf(jsonStruct)
	for i := range jsonStructReflectType.NumField() {
		tag, ok := jsonStructReflectType.Field(i).Tag.Lookup("json")
		// TODO: silently drops broken keys, not sure if that's okay yet
		if ok {
			jsonTags = append(jsonTags, tag)
			jsonTagPositions[tag] = i
		}
	}

	selectionName, err := Select(jsonTags, func(s string) string {
		return s
	}, reader)
	if err != nil {
		return "", 0, err
	}

	selectionIndex, ok := jsonTagPositions[selectionName]
	if !ok {
		return "", 0, fmt.Errorf("unexpected error: selection of json-tag did not return a valid key")
	}
	return selectionName, selectionIndex, nil
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
