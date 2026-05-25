package scenes

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/SaschaRunge/llm-local/internal/communication"
	"github.com/SaschaRunge/llm-local/internal/core"
	"github.com/SaschaRunge/llm-local/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var DefaultUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type mockRuntime struct {
	db      *database.Queries
	inputs  []string
	counter int
}

func (r *mockRuntime) Context() context.Context {
	return context.Background()
}

func (r *mockRuntime) CurrentScene() core.Scene {
	return &Chat{}
}

func (r *mockRuntime) DB() *database.Queries {
	return r.db
}

func (r *mockRuntime) GetInput() string {
	return ""
}

func (r *mockRuntime) Handle(input string) (core.Scene, error) {
	if input == "/exit" {
		return &Chat{}, nil
	}
	return nil, nil
}

/*
func (r *mockRuntime) GetInput() string {
	if len(r.inputs) == 0 {
		return ""
	}
	defer func() {
		r.counter += 1
	}()
	return r.inputs[r.counter%len(r.inputs)]
}*/

func TestSceneChat(t *testing.T) {
	dbQueries := setupChat(t,
		[]string{"Mexico", "Wunderland"},
		[]string{"Dickerchen", "Traube", "Wuerstchen", "Stubenhocker"})
	runtime := &mockRuntime{db: dbQueries, inputs: []string{"test1", "test2", "test3"}}

	chats, err := dbQueries.GetAllChats(runtime.Context())
	if err != nil {
		t.Fatalf("setup: failed fetch chats: %v", err)
	}

	authors, err := dbQueries.GetCharactersInChat(runtime.Context(), chats[0].ID)
	if err != nil {
		t.Fatalf("setup: failed fetch characters from chat %q: %v", chats[0].Name, err)
	}

	messagesByAuthor := map[uuid.UUID][]string{
		authors[0].ID: {"Dickerchen1", "Dickerchen2", "Dickerchen3", "Dickerchen4"},
		authors[2].ID: {"Wuerstchen1", "Wuerstchen2"},
		authors[3].ID: {"Stubenhocker1", "Stubenhocker2", "Stubenhocker3"},
	}

	orderedAuthors := []uuid.UUID{
		authors[0].ID,
		authors[2].ID,
		authors[3].ID,
	}

	for i := range authors[0].ID {
		for _, authorID := range orderedAuthors {
			messageByAuthor := messagesByAuthor[authorID]
			if i >= len(messageByAuthor) {
				continue
			}

			_, err := dbQueries.AddMessage(runtime.Context(), database.AddMessageParams{
				Content:  messageByAuthor[i],
				AuthorID: authorID,
				ChatID:   chats[0].ID,
				Role:     string(communication.RoleUser),
			})
			if err != nil {
				t.Fatalf("failed to add message %q to chat %q with err %v", messageByAuthor[i], chats[0].Name, err)
			}

			//t.Logf("added message %q\n", msg.ContentAnswer)
		}
	}

	sceneChat, err := NewSceneChat(runtime, chats[0]) //, Character{ID: uuid.UUID{}, Name: "CoolName"})
	if err != nil {
		t.Errorf("failed to load data into sceneChat: %v", err)
	}

	expected := []string{"Dickerchen1", "Wuerstchen1", "Stubenhocker1", "Dickerchen2", "Wuerstchen2", "Stubenhocker2", "Dickerchen3", "Stubenhocker3", "Dickerchen4"}
	if len(sceneChat.cachedMessages) != len(expected) {
		t.Errorf("fetched %d messages, expected %d:", len(sceneChat.cachedMessages), len(expected))
	}
	for i, msg := range sceneChat.cachedMessages {
		if msg.Message.Content != expected[i] {
			t.Errorf("mismatch at message %d. Got %q, expected %q.", i, expected[i], msg.Message.Content)
		}
	}

	runtime.inputs = []string{"gibberish", "more gibberish", "/exit"}
	for i := range len(runtime.inputs) {
		sceneTest, err := sceneChat.HandleRawInput(runtime.inputs[i])
		if err != nil {
			t.Errorf("unexpected error during sceneChat.Execute: %v", err)
		}
		if sceneTest.NextScene == nil {
			t.Errorf("sceneChat.Execute exited without returning a new scene")
		}
	}
}

func setupChat(t *testing.T, chatNames, characterNames []string) *database.Queries {
	db := loadDatabase(t)
	clearDatabase(t, db)
	dbQueries := database.New(db)

	chats := make([]database.Chat, len(chatNames))
	err := errors.New("")

	for i, cn := range chatNames {
		chats[i], err = dbQueries.AddChat(context.Background(), database.AddChatParams{
			Name:            cn,
			UserCharacterID: DefaultUserID,
		})
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

	return dbQueries
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
