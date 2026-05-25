package commands

import (
	"fmt"

	"github.com/SaschaRunge/llm-local/internal/cli/commands/selector"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
)

func selectCharacter(commandCtx core.CommandContext, dbQueries *database.Queries) (character, error) {
	availableCharacters, err := dbQueries.GetAvailableCharacters(commandCtx.Runtime.Context())
	if err != nil {
		return character{}, err
	}
	if len(availableCharacters) == 0 {
		return character{}, fmt.Errorf("no valid characters available")
	}

	selectedCharacter := selector.Select(
		availableCharacters,
		func(char character) string { return char.Name },
		commandCtx.Runtime.GetInput)

	if selectedCharacter.Name == "" {
		return character{}, nil
	}

	return selectedCharacter, nil
}
