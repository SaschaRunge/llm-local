-- name: AddChat :one
INSERT INTO chats(id, created_at, updated_at, name)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: GetChatByID :one
SELECT * FROM chats
WHERE id = $1;

-- name: GetChatsLikeName :many
SELECT * FROM chats
WHERE name LIKE $1;

-- name: GetAllChats :many
SELECT * FROM chats;

-- name: GetUserCharacterInChatByID :one
SELECT id, name, card FROM characters
WHERE characters.id = (
    SELECT user_character_id 
    FROM chats
    WHERE chats.id = $1
);

-- name: GetCardFromChatByID :one
SELECT card FROM chats
WHERE chats.id = $1;
