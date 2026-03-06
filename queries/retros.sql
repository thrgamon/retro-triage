-- name: CreateRetro :one
INSERT INTO retros (name)
VALUES ($1)
RETURNING *;

-- name: GetRetro :one
SELECT * FROM retros WHERE id = $1;

-- name: ListRetros :many
SELECT * FROM retros ORDER BY created_at DESC LIMIT 50;
