-- +migrate Up
ALTER TABLE `wordpress_gbps` ADD COLUMN `start_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `maps_url`;

-- +migrate Down
ALTER TABLE `wordpress_gbps` DROP COLUMN `start_date`;
