-- +goose Up
ALTER TABLE users
ADD COLUMN apikey VARCHAR(64) NOT NULL UNIQUE;

-- +goose Down
ALTER TABLE users
DROP COLUMN apikey;