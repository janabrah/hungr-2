-- +goose Up
-- +goose StatementBegin
ALTER TABLE recipes
    ADD COLUMN source TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE recipes
    DROP COLUMN source;
-- +goose StatementEnd
