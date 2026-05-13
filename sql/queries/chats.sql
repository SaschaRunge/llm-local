-- name: AddChat :one
INSERT INTO chats(id, created_at, updated_at, name)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: GetChatsLikeName :many
SELECT id FROM chats
WHERE name LIKE $1;