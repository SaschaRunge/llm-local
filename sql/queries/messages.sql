-- name: AddMessage :one
INSERT INTO messages(id, created_at, updated_at, reasoning, content, chat_id, author_id, role, idx)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    (SELECT COALESCE(MAX(idx), 0) + 1 FROM messages WHERE chat_id = $3)
)
RETURNING *;

-- name: GetMessagesByChatID :many
SELECT * FROM messages
WHERE chat_id = $1
ORDER BY idx ASC;

-- name: GetChatHistory :many
SELECT messages.id, 
messages.reasoning, 
messages.content,
messages.author_id,
messages.role,
characters.name AS author_name 
FROM messages
INNER JOIN characters
    ON messages.author_id = characters.id
WHERE messages.chat_id = $1
ORDER BY idx ASC;


-- name: DeleteMessages :exec
DELETE FROM messages;