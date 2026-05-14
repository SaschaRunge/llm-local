package scenes

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type runtime struct {
	db *database.Queries
}

func (r *runtime) Context() context.Context {
	return context.Background()
}
func (r *runtime) DB() *database.Queries {
	return r.db
}
func (r *runtime) ExecuteCommand(input string) (core.Scene, error) {
	return nil, nil
}
func (r *runtime) GetInput() string {
	return ""
}

func TestSceneChat(t *testing.T) {
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
		return
	}

	clearDatabase(t, db)
	dbQueries := database.New(db)

	chatNames := []string{"Mexico", "Wunderland"}
	chats := make([]database.Chat, len(chatNames))

	for i, cn := range chatNames {
		chats[i], err = dbQueries.AddChat(context.Background(), cn)
		if err != nil {
			t.Fatalf("setup: failed to add chat %q: %v", cn, err)
		}
	}

	characterNames := []string{"Dickerchen", "Traube", "Wuerstchen", "Stubenhocker"}
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

	err = sceneChat.loadData(&runtime{db: dbQueries})
	if err != nil {
		t.Errorf("failed to load data into sceneChat: %v", err)
	}
}

func clearDatabase(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE chats, messages, characters, chat_subscriptions RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}
}
