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