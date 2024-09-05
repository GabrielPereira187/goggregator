-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, feed_id, title, url, description, published_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)    
RETURNING *;

-- name: GetAllRecentPostsByUser :many
SELECT p.* FROM posts p INNER JOIN feeds f ON f.feed_id = p.feed WHERE f.user_id = $1 ORDER BY p.published_at DESC LIMIT $2;

