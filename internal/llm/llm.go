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
	model    string
	provider *llamacpp.Provider
}

// TODO: create client with empty history
func NewClient(model string) (*Client, error) {
	client := &Client{}
	if err := client.load(model); err != nil {
		return &Client{}, err
	}

	return client, nil
}

func (c *Client) load(model string) error {
	provider, err := llamacpp.New()
	if err != nil {
		return err
	}

	c.provider = provider
	c.model = model

	return nil
}

func (c *Client) GenerateAnswer(ctx context.Context, messageHistory []communication.Message) (string, error) {
	anyllmMessageHistory, err := translateHistory(messageHistory)
	if err != nil {
		return "", nil
	}

	response, err := c.provider.Completion(ctx, anyllm.CompletionParams{
		Model:    c.model,
		Messages: anyllmMessageHistory,
	})

	if err != nil || len(response.Choices) == 0 {
		return "", err
	}

	return response.Choices[0].Message.ContentString(), nil
}

func (c *Client) StreamAnswer(input string) string {
	return ""
}

func translate(message communication.Message) (anyllm.Message, error) {
	if !roleIsValid(message.Role) {
		return anyllm.Message{}, fmt.Errorf("message role %q is not a valid role", message.Role)
	}
	return anyllm.Message{
		Name:    message.Name,
		Role:    string(message.Role),
		Content: message.Content,
		//Reasoning: message.Reasoning, //TODO: might need specific mapping, not sure yet
	}, nil
}

func translateHistory(messageHistory []communication.Message) ([]anyllm.Message, error) {
	anyllmMessageHistory := []anyllm.Message{}
	for _, message := range messageHistory {
		anyllmMessage, err := translate(message)
		if err != nil {
			return []anyllm.Message{}, err
		}

		anyllmMessageHistory = append(anyllmMessageHistory, anyllmMessage)
	}

	return anyllmMessageHistory, nil
}

func roleIsValid(role communication.Role) bool {
	switch role {
	case anyllm.RoleSystem, anyllm.RoleAssistant, anyllm.RoleUser:
		return true
	default:
		return false
	}
}

func loadSystemPrompt() (string, error) {
	systemPrompt, err := os.ReadFile(pathToSystemPrompt)
	if err != nil {
		return "", err
	}

	return string(systemPrompt), err
}
