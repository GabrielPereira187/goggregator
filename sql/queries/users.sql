-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, apikey)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users where apikey = $1;