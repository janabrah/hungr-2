-- +goose Up

-- Create trigger function for auto-updating updated_at
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- Add updated_at to recipes (already has created_at)
ALTER TABLE recipes ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();

-- Add timestamps to files
ALTER TABLE files ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE files ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();

-- Add timestamps to tags
ALTER TABLE tags ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE tags ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();

-- Add timestamps to recipe_tags
ALTER TABLE recipe_tags ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE recipe_tags ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();

-- Create triggers for auto-updating updated_at
CREATE TRIGGER update_recipes_updated_at
    BEFORE UPDATE ON recipes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_files_updated_at
    BEFORE UPDATE ON files
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tags_updated_at
    BEFORE UPDATE ON tags
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_recipe_tags_updated_at
    BEFORE UPDATE ON recipe_tags
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose Down

-- Drop triggers
DROP TRIGGER IF EXISTS update_recipe_tags_updated_at ON recipe_tags;
DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;
DROP TRIGGER IF EXISTS update_files_updated_at ON files;
DROP TRIGGER IF EXISTS update_recipes_updated_at ON recipes;

-- Drop columns from recipe_tags
ALTER TABLE recipe_tags DROP COLUMN updated_at;
ALTER TABLE recipe_tags DROP COLUMN created_at;

-- Drop columns from tags
ALTER TABLE tags DROP COLUMN updated_at;
ALTER TABLE tags DROP COLUMN created_at;

-- Drop columns from files
ALTER TABLE files DROP COLUMN updated_at;
ALTER TABLE files DROP COLUMN created_at;

-- Drop column from recipes
ALTER TABLE recipes DROP COLUMN updated_at;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();
