-- +goose Up
-- +goose StatementBegin
CREATE TYPE ingredient_unit AS ENUM ('ml', 'mg', 'count');

CREATE TABLE ingredient_names (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE recipe_steps (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_uuid UUID NOT NULL REFERENCES recipes(uuid) ON DELETE CASCADE,
    step_number INT NOT NULL,
    instructions TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_recipe_steps_recipe_uuid ON recipe_steps(recipe_uuid);
CREATE UNIQUE INDEX idx_recipe_steps_recipe_step ON recipe_steps(recipe_uuid, step_number);

CREATE TABLE step_ingredients (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_step_uuid UUID NOT NULL REFERENCES recipe_steps(uuid) ON DELETE CASCADE,
    ingredient_name_uuid UUID NOT NULL REFERENCES ingredient_names(uuid),
    ingredient_type ingredient_unit NOT NULL,
    quantity DECIMAL NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_step_ingredients_recipe_step ON step_ingredients(recipe_step_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS step_ingredients;
DROP TABLE IF EXISTS recipe_steps;
DROP TABLE IF EXISTS ingredient_names;
DROP TYPE IF EXISTS ingredient_unit;
-- +goose StatementEnd
