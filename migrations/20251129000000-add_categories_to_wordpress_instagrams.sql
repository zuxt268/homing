-- +migrate Up
ALTER TABLE `wordpress_instagrams` ADD COLUMN `categories` varchar(255) NOT NULL DEFAULT "" AFTER `memo`;

-- +migrate Down
ALTER TABLE `wordpress_instagrams` DROP COLUMN `categories`;