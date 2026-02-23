-- +migrate Up
CREATE TABLE IF NOT EXISTS `wordpress_gbps` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `memo` text,
    `wordpress_domain` varchar(255) NOT NULL,
    `business_name` varchar(255) NOT NULL,
    `business_title` varchar(255) NOT NULL,
    `maps_url` varchar(255) DEFAULT '',
    `status` int NOT NULL DEFAULT '0',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +migrate Down
DROP TABLE `wordpress_gbps`;
