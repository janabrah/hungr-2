-- +goose Up
ALTER TABLE recipes DROP COLUMN tag_string;
ALTER TABLE recipe_tags ADD COLUMN id SERIAL;

-- +goose Down
ALTER TABLE recipe_tags DROP COLUMN id;
ALTER TABLE recipes ADD COLUMN tag_string TEXT;
