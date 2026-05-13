-- name: SubscribeToChat :one
INSERT INTO chat_subscriptions(id, chat_id, character_id)
VALUES(
    gen_random_uuid(),
    $1,
    $2
)
RETURNING *;

-- name: GetCharactersLikeName :many
SELECT id FROM characters 
WHERE name LIKE $1;