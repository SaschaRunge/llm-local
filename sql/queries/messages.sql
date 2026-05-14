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

-- name: GetMessagesByChatID :many
SELECT * FROM messages
WHERE chat_id = $1;

-- name: GetChatHistory :many
SELECT messages.id, 
messages.content_thoughts, 
messages.content_answer,
messages.author_id, 
characters.name AS author_name 
FROM messages
INNER JOIN characters
    ON messages.author_id = characters.id
WHERE messages.chat_id = $1;


-- name: DeleteMessages :exec
DELETE FROM messages;