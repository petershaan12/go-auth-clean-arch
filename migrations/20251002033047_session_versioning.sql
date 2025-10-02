-- +goose Up
ALTER TABLE users ADD COLUMN session_version INT NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE users DROP COLUMN session_version;