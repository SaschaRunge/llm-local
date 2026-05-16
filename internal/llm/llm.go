package llm

import (
	//anyllm "github.com/mozilla-ai/any-llm-go"
	"github.com/mozilla-ai/any-llm-go/providers/llamacpp"
)

type Client struct {
	provider *llamacpp.Provider
}

func (c *Client) load() error {
	provider, err := llamacpp.New()
	if err != nil {
		return err
	}

	c.provider = provider

	return nil
}
