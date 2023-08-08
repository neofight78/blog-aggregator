-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (UUID_TO_BIN(sqlc.arg(id)), ?, ?, ?, ?, ?, ?, UUID_TO_BIN(sqlc.arg(feed_id)));

-- name: GetPostsByUser :many
SELECT p.id, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at, p.feed_id
FROM posts p JOIN feed_follows ON p.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = UUID_TO_BIN(sqlc.arg(user_id)) LIMIT ?
