-- +goose Up
CREATE TABLE users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- Add foreign key constraint to recipes table
ALTER TABLE recipes
ADD CONSTRAINT fk_recipes_user_uuid
FOREIGN KEY (user_uuid) REFERENCES users(uuid);

-- +goose Down
ALTER TABLE recipes DROP CONSTRAINT IF EXISTS fk_recipes_user_uuid;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
