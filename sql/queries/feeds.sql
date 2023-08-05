-- name: CreateFeed :exec
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (UUID_TO_BIN(sqlc.arg(id)), ?, ?, ?, ?, UUID_TO_BIN(sqlc.arg(user_id)));

-- name: ListFeeds :many
SELECT id, created_at, updated_at, name, url, user_id FROM feeds;
