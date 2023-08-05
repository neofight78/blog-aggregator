-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, name, api_key)
VALUES (UUID_TO_BIN(sqlc.arg(id)), ?, ?, ?, sha2(UUID(), 256));

-- name: GetUser :one
SELECT id, created_at, updated_at, name, api_key FROM users WHERE api_key = ?
