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

-- name: AddVariant :one
INSERT INTO messages(id, created_at, updated_at, reasoning, content, chat_id, author_id, role, idx, parent_msg_id, sequence_idx)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    (SELECT idx FROM messages WHERE id = $6),
    $6,
    (SELECT COALESCE(MAX(sequence_idx), 0) + 1 FROM messages WHERE parent_msg_id = $6)
)
RETURNING *;

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

-- name: GetVariants :many
SELECT messages.id, 
messages.reasoning, 
messages.content,
messages.author_id,
messages.role,
characters.name AS author_name 
FROM messages
INNER JOIN characters
    ON messages.author_id = characters.id
WHERE messages.parent_msg_id = $1
ORDER BY sequence_idx ASC;


-- name: DeleteMessages :exec
DELETE FROM messages;

-- name: ReplaceMessage :exec
UPDATE messages
SET 
    updated_at = NOW(),
    reasoning = $2,
    content = $3
WHERE id = $1;

