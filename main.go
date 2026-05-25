package main

import (
	_ "github.com/lib/pq"

	_ "bufio"

	_ "context"

	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SaschaRunge/llm-local/internal/cli"
	_ "github.com/SaschaRunge/llm-local/internal/cli"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"

	_ "github.com/google/uuid"

	"github.com/joho/godotenv"

	anyllm "github.com/mozilla-ai/any-llm-go"

	_ "github.com/mozilla-ai/any-llm-go/providers/llamacpp"
)

const pathToSystemPrompt = "./system_prompt.md"
const colorGrey = "\033[2m"
const colorReset = "\033[0m"

func main() {
	//ctx := context.Background()

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("unable to open database: %s", err)
		return
	}

	dbQueries := database.New(db)

	ui := cli.New(&core.Store{
		DB:        db,
		DBQueries: dbQueries})
	ui.Run()

	/*
		//ui := cli.New(dbQueries)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		//input := scanner.Text()
		/*

			_, err = ui.ExecuteCommand(input)
			if err != nil {
				fmt.Printf("Unable to run command %s: %s", input, err)
			} */
	/*
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

		//authorIDSystem := uuid.New()
		_ = uuid.New()

		charIDs, err := dbQueries.GetCharactersLikeName(context.Background(), "Toertchen")
		chatIDs, err := dbQueries.GetChatsLikeName(context.Background(), "%test%")

		//fmt.Printf("CHAT ID: %v", chatIDs)

		_, err = dbQueries.AddMessage(context.Background(), database.AddMessageParams{
			ContentAnswer: string(systemPrompt),
			ChatID:        chatIDs[0].ID,
			AuthorID:      charIDs[0],
		})

		if err != nil {
			fmt.Println(err)
		}

		char1, _ := dbQueries.AddCharacter(context.Background(), database.AddCharacterParams{
			Name: "Kathi",
		})
		char2, _ := dbQueries.AddCharacter(context.Background(), database.AddCharacterParams{
			Name: "Sascha",
		})
		dbQueries.AddCharacter(context.Background(), database.AddCharacterParams{
			Name: "Peter",
		})

		chat, _ := dbQueries.AddChat(context.Background(), "UnserChat")
		dbQueries.SubscribeToChat(context.Background(), database.SubscribeToChatParams{
			ChatID:      chat.ID,
			CharacterID: char1.ID,
		})
		dbQueries.SubscribeToChat(context.Background(), database.SubscribeToChatParams{
			ChatID:      chat.ID,
			CharacterID: char2.ID,
		})

		//chats, _ := dbQueries.GetChatsLikeName(context.Background(), "%Unser%")

		characters, _ := dbQueries.GetCharactersInChat(context.Background(), chat.ID)

		for _, char := range characters {
			fmt.Printf("Character: %s\n", char.Name)
		}

		if true {
			return
		}

		for {
			fmt.Print("User:")
			scanner.Scan()
			input := scanner.Text()
			message := anyllm.Message{Role: anyllm.RoleUser, Content: input}
			messages = append(messages, message)
			users, _ := dbQueries.GetCharactersLikeName(context.Background(), "%user%")
			dbQueries.AddMessage(context.Background(), database.AddMessageParams{
				ContentAnswer: string(input),
				ChatID:        chatIDs[0].ID,
				AuthorID:      users[0],
			})
			chunks, errs := provider.CompletionStream(ctx, anyllm.CompletionParams{
				Model:    "qwen3.6-27b",
				Messages: messages,
			})

			answer, _, _ := handleStream(chunks, errs)
			messages = append(messages, anyllm.Message{Role: anyllm.RoleAssistant, Content: answer})
			dbQueries.AddMessage(context.Background(), database.AddMessageParams{
				ContentAnswer: string(answer),
				ChatID:        chatIDs[0].ID,
				AuthorID:      charIDs[0],
			})

		}
	*/
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
