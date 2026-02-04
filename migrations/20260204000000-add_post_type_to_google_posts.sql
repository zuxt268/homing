-- +migrate Up
ALTER TABLE google_posts ADD COLUMN post_type VARCHAR(16) DEFAULT 'photo';

-- +migrate Down
ALTER TABLE google_posts DROP COLUMN post_type;