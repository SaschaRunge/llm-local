-- name: SubscribeToChat :exec
INSERT INTO chat_subscriptions(id, chat_id, character_id)
VALUES(
    gen_random_uuid(),
    $1,
    $2
);

-- name: GetCharactersInChat :many
SELECT 
characters.id, 
characters.name,
characters.card
FROM chat_subscriptions 
INNER JOIN characters
    ON chat_subscriptions.character_id = characters.id
WHERE 
    chat_subscriptions.chat_id = $1 AND
    characters.deleted_at IS NULL;;

