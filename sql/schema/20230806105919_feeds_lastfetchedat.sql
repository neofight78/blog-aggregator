-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds
    ADD COLUMN last_fetched_at DATETIME NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds
    DROP COLUMN last_fetched_at;
-- +goose StatementEnd
