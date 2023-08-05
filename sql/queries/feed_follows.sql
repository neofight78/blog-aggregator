-- name: CreateFeedFollow :exec
INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
VALUES (UUID_TO_BIN(sqlc.arg(id)), ?, ?, UUID_TO_BIN(sqlc.arg(feed_id)), UUID_TO_BIN(sqlc.arg(user_id)));

-- name: GetFeedFollow :one
SELECT id, created_at, updated_at, feed_id, user_id FROM feed_follows WHERE id = UUID_TO_BIN(sqlc.arg(id));

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id = UUID_TO_BIN(sqlc.arg(id));

-- name: ListFeedFollows :many
SELECT id, created_at, updated_at, feed_id, user_id FROM feed_follows WHERE user_id = UUID_TO_BIN(sqlc.arg(id))
