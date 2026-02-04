
-- +migrate Up
ALTER TABLE `google_posts` DROP COLUMN `google_business_url`;

-- +migrate Down
ALTER TABLE `google_posts` ADD COLUMN `google_business_url` varchar(500) NOT NULL DEFAULT '' AFTER `id`;