package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/SaschaRunge/llm-local/internal/communication"

	anyllm "github.com/mozilla-ai/any-llm-go"
	"github.com/mozilla-ai/any-llm-go/providers/llamacpp"
)

const pathToSystemPrompt = "./system_prompt.md"

type Client struct {
	messageHistory []anyllm.Message
	model          string
	provider       *llamacpp.Provider
}

// TODO: create client with empty history
func NewClient(messageHistory []communication.Message) (*Client, error) {
	client := &Client{}
	if err := client.load(messageHistory); err != nil {
		return &Client{}, err
	}

	return client, nil
}

func (c *Client) load(messageHistory []communication.Message) error {
	provider, err := llamacpp.New()
	if err != nil {
		return err
	}

	c.provider = provider

	if len(messageHistory) == 0 {
		systemPrompt, err := loadSystemPrompt()
		if err != nil {
			return err
		}
		c.messageHistory = append(c.messageHistory, anyllm.Message{
			Role:    anyllm.RoleSystem,
			Content: systemPrompt,
		})
	}

	for i, message := range messageHistory {
		if i == 0 {
			if message.Role != anyllm.RoleSystem {
				return fmt.Errorf("unable to load llm client: first message in history is not system prompt")
			}
		}
		c.messageHistory = append(c.messageHistory, translate(message))
	}

	return nil
}

func (c *Client) GenerateAnswer(input string) (string, error) {
	response, err := c.provider.Completion(context.Background(), anyllm.CompletionParams{
		Model: "qwen3.6-27b",
		Messages: append(c.messageHistory, anyllm.Message{
			Content: input,
			Role:    anyllm.RoleUser,
		}),
	})

	if err != nil || len(response.Choices) == 0 {
		return "", nil
	}

	return response.Choices[0].Message.ContentString(), nil
}

func (c *Client) StreamAnswer(input string) string {
	return ""
}

func translate(message communication.Message) anyllm.Message {
	return anyllm.Message{
		Name:    message.Name,
		Role:    string(message.Role),
		Content: message.Content,
		//Reasoning: message.Reasoning, //TODO: might need specific mapping, not sure yet
	}
}

func loadSystemPrompt() (string, error) {
	systemPrompt, err := os.ReadFile(pathToSystemPrompt)
	if err != nil {
		return "", err
	}

	return string(systemPrompt), err
}
