-- +goose Up
ALTER TABLE users ADD COLUMN last_seen TIMESTAMPTZ DEFAULT NOW();

-- +goose Down
ALTER TABLE users DROP COLUMN last_seen;
