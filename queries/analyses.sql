-- name: UpsertAnalysis :one
INSERT INTO analyses (retro_id, result)
VALUES ($1, $2)
ON CONFLICT (retro_id) DO UPDATE SET result = EXCLUDED.result, created_at = now()
RETURNING *;

-- name: GetAnalysisByRetroID :one
SELECT * FROM analyses WHERE retro_id = $1;
