-- +migrate Up
ALTER TABLE google_posts DROP COLUMN media_format;

-- +migrate Down
ALTER TABLE google_posts ADD COLUMN media_format VARCHAR(64);