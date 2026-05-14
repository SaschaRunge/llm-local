package llm

import (
	//anyllm "github.com/mozilla-ai/any-llm-go"
	"github.com/mozilla-ai/any-llm-go/providers/llamacpp"
)

type LLM struct {
	provider *llamacpp.Provider
}

func (l *LLM) load() error {
	provider, err := llamacpp.New()
	if err != nil {
		return err
	}

	l.provider = provider

	return nil
}
