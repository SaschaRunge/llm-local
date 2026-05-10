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

const pathToSystemPrompt = "./system_prompt.md"
const colorGrey = "\033[2m"
const colorReset = "\033[0m"

func main() {
	ctx := context.Background()

	provider, err := llamacpp.New()
	if err != nil {
		log.Fatal(err)
	}

	systemPrompt, err := os.ReadFile(pathToSystemPrompt)
	if err != nil {
		log.Fatal(err)
	}

	messages := []anyllm.Message{
		{Role: anyllm.RoleSystem, Content: string(systemPrompt)},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("User:")
		scanner.Scan()
		input := scanner.Text()
		message := anyllm.Message{Role: anyllm.RoleUser, Content: input}
		messages = append(messages, message)
		chunks, errs := provider.CompletionStream(ctx, anyllm.CompletionParams{
			Model:    "qwen3.6-27b",
			Messages: messages,
		})

		answer, _, _ := handleStream(chunks, errs)
		messages = append(messages, anyllm.Message{Role: anyllm.RoleAssistant, Content: answer})

	}
}

// this is not proper: very short thinks will break, </think> is not greyed and parts of the thoughts may leak in to the answer
// maybe adjust buffer size by content or just buffer if < is present
func handleStream(chunks <-chan anyllm.ChatCompletionChunk, errs <-chan error) (answer, thoughts string, err error) {
	var answerContent strings.Builder
	var thoughtsContent strings.Builder

	isThinking := true

	fmt.Print(colorGrey)
	for chunk := range chunks {
		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if isThinking {
				thoughtsContent.WriteString(content)
				if len(thoughtsContent.String()) > 20 && strings.Contains(string(thoughtsContent.String()[len(thoughtsContent.String())-20:]), "</think>") {
					fmt.Print(colorReset)
					isThinking = false
				}
			} else {
				answerContent.WriteString(content)
			}
			fmt.Print(content)
			os.Stdout.Sync()
		}
	}

	for err = range errs {
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println()

	return answerContent.String(), thoughtsContent.String(), nil
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
