package parser

import (
	"strings"
)

type preprocessor func(string) (string, string)

func Default(rawInput string)

func SelectPreprocessorByPrefix(rawInput string) (pp preprocessor, ok bool) {
	for prefix, pp := range preprocessors() {
		if strings.HasPrefix(rawInput, prefix) {
			return pp, true
		}
	}
	return nil, false
}

func preprocessors() map[string]preprocessor {
	return map[string]preprocessor{
		"/": preprocessSlash,
		"@": preprocessMention,
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

func preprocessMention(rawInput string) (cmd string, rawArgs string) {
	cmd = "@"

	//fix, throws name outa window
	parts := strings.SplitN(rawInput, " ", 2)
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
