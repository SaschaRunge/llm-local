-- name: AddChat :one
INSERT INTO chats(id, created_at, updated_at, name, card, user_character_id)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
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

-- name: GetAvailableChats :many
SELECT id, name, card FROM chats
WHERE
    deleted_at IS NULL;

-- name: UpdateChatCard :one
UPDATE chats
SET 
    updated_at = NOW(),
    card = $2
WHERE id = $1
RETURNING id;
