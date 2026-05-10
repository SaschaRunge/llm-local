package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	anyllm "github.com/mozilla-ai/any-llm-go"
	"github.com/mozilla-ai/any-llm-go/providers/llamacpp"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	provider, err := llamacpp.New()
	if err != nil {
		log.Fatal(err)
	}

	messages := []anyllm.Message{
		{Role: anyllm.RoleSystem, Content: "You are the satan."},
	}

	for {
		scanner.Scan()
		input := scanner.Text()
		message := anyllm.Message{Role: anyllm.RoleUser, Content: input}
		messages = append(messages, message)
		response, err := provider.Completion(ctx, anyllm.CompletionParams{
			Model:    "qwen3.6-27b",
			Messages: messages,
		})
		if err != nil {
			log.Fatal(err)
		}
		answer, thoughts := extractThinking(response.Choices[0].Message.ContentString())
		fmt.Println("_________________________")
		fmt.Printf("%s\n", thoughts)
		fmt.Println("_________________________")
		fmt.Printf("%s\n", answer)
		messages = append(messages, anyllm.Message{Role: anyllm.RoleAssistant, Content: answer})

		/*fmt.Println("DEBUG:")
		for _, msg := range messages {

			fmt.Println(msg.ContentString())
		}
		fmt.Println("DEBUG END")*/
	}
}

func extractThinking(content string) (answer, thoughts string) {
	parts := strings.Split(content, "</think>")

	if len(parts) > 1 {
		answer := parts[1]
		thoughts := strings.Trim(parts[0], "<think>")
		thoughts = strings.Trim(thoughts, "</think>")
		thoughts = strings.TrimSpace(thoughts)

		return answer, thoughts
	}

	return parts[0], ""
}
