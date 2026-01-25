-- +goose Up
ALTER TABLE recipes ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE recipes DROP COLUMN is_public;
