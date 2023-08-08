-- name: CreateFeed :exec
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (UUID_TO_BIN(sqlc.arg(id)), ?, ?, ?, ?, UUID_TO_BIN(sqlc.arg(user_id)));

-- name: ListFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds ORDER BY last_fetched_at LIMIT ?;

-- name: MarkFeedFetched :exec
UPDATE feeds SET updated_at = ?, last_fetched_at = ? WHERE id = UUID_TO_BIN(sqlc.arg(id));
