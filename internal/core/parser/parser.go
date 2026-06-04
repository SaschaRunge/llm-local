package parser

import (
	"fmt"
	"strings"
)

type preprocessor func(string) (string, string)

var preprocessors = map[string]preprocessor{
	"/": preprocessSlash,
	"@": createMentionPreprocessor("@"),
}

func SelectPreprocessorByPrefix(rawInput string) (pp preprocessor, ok bool) {
	for prefix, pp := range preprocessors {
		if strings.HasPrefix(rawInput, prefix) {
			return pp, true
		}
	}
	return nil, false
}

func createMentionPreprocessor(prefix string) preprocessor {
	return func(rawInput string) (cmd string, rawArgs string) {
		rawArgs = strings.TrimPrefix(rawInput, prefix)
		return prefix, rawArgs
	}
}

func preprocessSlash(rawInput string) (cmd string, rawArgs string) {
	parts := strings.SplitN(rawInput, " ", 2)

	cmd = parts[0]
	if len(parts) >= 2 {
		rawArgs = parts[1]
	}

	return cmd, rawArgs
}

func ParseArgs(rawArgs string) []string {
	if rawArgs == "" {
		return nil
	}
	return strings.Fields(rawArgs)
}

// TODO: no longer belongs to parser and should be part of llm. temporary workaround
func ExtractReasoning(input, closingTag string) (content, reasoning string) {
	//openingTag := fmt.Sprintf("<%s>", tagName)

	iOpening := 0 //strings.Index(input, openingTag)
	iClosing := strings.LastIndex(input, closingTag)

	if iOpening < 0 || iClosing < 0 {
		return strings.TrimSpace(input), ""
	}

	content = fmt.Sprintf("%s%s", input[:iOpening], input[iClosing+len(closingTag):])
	reasoning = fmt.Sprintf("%s", input[0:iClosing]) //input[iOpening+len(openingTag):iClosing])

	return strings.TrimSpace(content), strings.TrimSpace(reasoning)
}
