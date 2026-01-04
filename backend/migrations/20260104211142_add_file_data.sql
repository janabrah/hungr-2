-- +goose Up
ALTER TABLE files ADD COLUMN data BYTEA;
ALTER TABLE files ADD COLUMN content_type TEXT DEFAULT 'image/jpeg';

-- +goose Down
ALTER TABLE files DROP COLUMN content_type;
ALTER TABLE files DROP COLUMN data;
