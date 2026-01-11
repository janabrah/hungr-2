-- +goose Up
CREATE TABLE user_connections (
    source_user_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    target_user_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    PRIMARY KEY (source_user_uuid, target_user_uuid),
    CONSTRAINT no_self_connection CHECK (source_user_uuid != target_user_uuid)
);

CREATE INDEX idx_user_connections_source ON user_connections(source_user_uuid);
CREATE INDEX idx_user_connections_target ON user_connections(target_user_uuid);

-- +goose Down
DROP INDEX IF EXISTS idx_user_connections_target;
DROP INDEX IF EXISTS idx_user_connections_source;
DROP TABLE IF EXISTS user_connections;
