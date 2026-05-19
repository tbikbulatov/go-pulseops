-- +goose Up
ALTER TABLE incidents
ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE incidents
DROP COLUMN IF EXISTS version;
