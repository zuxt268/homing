
-- +migrate Up
ALTER TABLE `business_instagrams` ADD COLUMN `maps_url` varchar(512) NOT NULL DEFAULT '' AFTER `business_title`;

-- +migrate Down
ALTER TABLE `business_instagrams` DROP COLUMN `maps_url`;