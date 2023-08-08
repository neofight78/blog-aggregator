-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts
(
    id BINARY(16) NOT NULL PRIMARY KEY,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(1023) NOT NULL,
    published_at DATETIME NOT NULL,
    feed_id BINARY(16) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
