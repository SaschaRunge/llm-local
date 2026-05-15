package scenes

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type mockRuntime struct {
	db *database.Queries
}

func (r *mockRuntime) Context() context.Context {
	return context.Background()
}
func (r *mockRuntime) DB() *database.Queries {
	return r.db
}
func (r *mockRuntime) ExecuteCommand(input string) (core.Scene, error) {
	return nil, nil
}
func (r *mockRuntime) GetInput() string {
	return ""
}

func TestSceneChat(t *testing.T) {
	dbQueries, sceneChat := setup(t,
		[]string{"Mexico", "Wunderland"},
		[]string{"Dickerchen", "Traube", "Wuerstchen", "Stubenhocker"})
	runtime := &mockRuntime{db: dbQueries}

	authors, err := dbQueries.GetCharactersInChat(runtime.Context(), sceneChat.ID)
	if err != nil {
		t.Fatalf("setup: failed fetch characters from chat %q: %v", sceneChat.Name, err)
	}

	messagesByAuthor := map[uuid.UUID][]string{
		authors[0].ID: {"Dickerchen1", "Dickerchen2", "Dickerchen3", "Dickerchen4"},
		authors[2].ID: {"Wuerstchen1", "Wuerstchen2"},
		authors[3].ID: {"Stubenhocker1", "Stubenhocker2", "Stubenhocker3"},
	}

	for i := range 4 {
		for k, v := range messagesByAuthor {
			if i >= len(v) {
				continue
			}

			_, err := dbQueries.AddMessage(runtime.Context(), database.AddMessageParams{
				ContentAnswer: v[i],
				AuthorID:      k,
				ChatID:        sceneChat.ID,
			})
			if err != nil {
				t.Fatalf("failed to add message %q to chat %q with err %v", v[i], sceneChat.Name, err)
			}

			//t.Logf("added message %q\n", msg.ContentAnswer)
		}
	}

	err = sceneChat.loadData(runtime)
	if err != nil {
		t.Errorf("failed to load data into sceneChat: %v", err)
	}

	/*for _, msg := range sceneChat.messages {
		t.Logf("messages: %q\n", msg.contentAnswer)
	}*/

	expected := []string{"Dickerchen1", "Wuerstchen1", "Stubenhocker1", "Dickerchen2", "Wuerstchen2", "Stubenhocker2", "Dickerchen3", "Stubenhocker3", "Dickerchen4"}
	if len(sceneChat.messages) != len(expected) {
		t.Errorf("fetched %d messages, expected %d:", len(sceneChat.messages), len(expected))
	}
	for i, msg := range sceneChat.messages {
		if msg.contentAnswer != expected[i] {
			t.Errorf("mismatch at message %d. Got %q, expected %q.", i, expected[i], msg.contentAnswer)
		}
	}
}

func setup(t *testing.T, chatNames, characterNames []string) (*database.Queries, SceneChat) {
	db := loadDatabase(t)
	clearDatabase(t, db)
	dbQueries := database.New(db)

	chats := make([]database.Chat, len(chatNames))
	err := errors.New("")

	for i, cn := range chatNames {
		chats[i], err = dbQueries.AddChat(context.Background(), cn)
		if err != nil {
			t.Fatalf("setup: failed to add chat %q: %v", cn, err)
		}
	}

	characters := make([]database.Character, len(characterNames))

	for i, cn := range characterNames {
		characters[i], err = dbQueries.AddCharacter(context.Background(), database.AddCharacterParams{Name: cn})
		if err != nil {
			t.Fatalf("setup: failed to add character %q: %v", cn, err)
		}
	}

	for _, char := range characters {
		err = dbQueries.SubscribeToChat(context.Background(), database.SubscribeToChatParams{
			CharacterID: char.ID,
			ChatID:      chats[0].ID})
		if err != nil {
			t.Fatalf("setup: failed to subscribe character %q to chat %q: %v", char.Name, chats[0].Name, err)
		}
	}

	sceneChat := SceneChat{
		ID:   chats[0].ID,
		Name: chats[0].Name,
	}

	return dbQueries, sceneChat
}

func loadDatabase(t *testing.T) *sql.DB {
	home, _ := os.UserHomeDir()
	envFile := filepath.Join(home, "workspace/github.com/SaschaRunge/Go/llm-local/.env")
	godotenv.Load(envFile)

	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		t.Fatal("unable to load db url from .env")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("unable to open database: %s", err)
	}
	return db
}

func clearDatabase(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE chats, messages, characters, chat_subscriptions RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}
}
