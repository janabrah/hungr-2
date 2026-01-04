CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS recipes (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    filename TEXT NOT NULL,
    user_uuid UUID NOT NULL,
    tag_string TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS files (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url TEXT NOT NULL,
    image BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS file_recipes (
    file_uuid UUID REFERENCES files(uuid),
    recipe_uuid UUID REFERENCES recipes(uuid),
    page_number INTEGER DEFAULT 0,
    PRIMARY KEY (file_uuid, recipe_uuid)
);

CREATE TABLE IF NOT EXISTS tags (
    uuid UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS recipe_tags (
    recipe_uuid UUID REFERENCES recipes(uuid),
    tag_uuid UUID REFERENCES tags(uuid),
    PRIMARY KEY (recipe_uuid, tag_uuid)
);

CREATE INDEX IF NOT EXISTS idx_recipes_user_uuid ON recipes(user_uuid);
CREATE INDEX IF NOT EXISTS idx_file_recipes_recipe_uuid ON file_recipes(recipe_uuid);
CREATE INDEX IF NOT EXISTS idx_recipe_tags_recipe_uuid ON recipe_tags(recipe_uuid);
