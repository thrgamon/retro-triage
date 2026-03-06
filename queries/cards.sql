-- name: CreateCard :one
INSERT INTO cards (retro_id, column_type, content, author_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCardsByRetroID :many
SELECT * FROM cards WHERE retro_id = $1 ORDER BY created_at ASC;

-- name: DeleteCard :exec
DELETE FROM cards WHERE id = $1 AND retro_id = $2;
