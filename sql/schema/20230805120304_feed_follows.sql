-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows
(
    id BINARY(16) NOT NULL PRIMARY KEY,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    feed_id BINARY(16) NOT NULL,
    user_id BINARY(16) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;
-- +goose StatementEnd
