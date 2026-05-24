package parser

import (
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
