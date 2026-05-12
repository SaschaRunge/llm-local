-- name: AddMessage :one
INSERT INTO messages(id, created_at, updated_at, content_thoughts, content_answer, chat_id, author_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;