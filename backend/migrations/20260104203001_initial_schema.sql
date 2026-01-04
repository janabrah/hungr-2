-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE recipes (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    user_uuid UUID NOT NULL,
    tag_string TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE files (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_uuid UUID NOT NULL REFERENCES recipes(uuid),
    url TEXT NOT NULL,
    page_number INTEGER DEFAULT 0,
    image BOOLEAN DEFAULT true
);

CREATE TABLE tags (
    uuid UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE recipe_tags (
    recipe_uuid UUID REFERENCES recipes(uuid),
    tag_uuid UUID REFERENCES tags(uuid),
    PRIMARY KEY (recipe_uuid, tag_uuid)
);

CREATE INDEX idx_recipes_user_uuid ON recipes(user_uuid);
CREATE INDEX idx_files_recipe_uuid ON files(recipe_uuid);
CREATE INDEX idx_recipe_tags_recipe_uuid ON recipe_tags(recipe_uuid);

-- +goose Down
DROP INDEX IF EXISTS idx_recipe_tags_recipe_uuid;
DROP INDEX IF EXISTS idx_files_recipe_uuid;
DROP INDEX IF EXISTS idx_recipes_user_uuid;
DROP TABLE IF EXISTS recipe_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS recipes;
