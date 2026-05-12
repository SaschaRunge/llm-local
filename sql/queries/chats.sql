-- name: AddChat :one
INSERT INTO chats(id, created_at, updated_at, name)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: GetChatByName :many
SELECT * FROM chats
WHERE name = $1;