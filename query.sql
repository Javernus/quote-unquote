-- name: Insert :one
INSERT INTO quote (id, message, person, created_at, updated_at, ip)
VALUES ($1, $2, $3, $4, $4, $5)
RETURNING *;

-- name: FindAll :many
SELECT *
FROM quote
ORDER BY created_at DESC
LIMIT $1;

-- name: Count :one
SELECT COUNT(*) FROM quote;
