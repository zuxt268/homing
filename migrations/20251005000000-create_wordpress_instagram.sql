
-- +migrate Up
CREATE TABLE IF NOT EXISTS `wordpress_instagrams` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `wordpress_domain` varchar(255) NOT NULL,
    `wordpress_site_title` varchar(255) NOT NULL,
    `instagram_id` varchar(255) NOT NULL,
    `instagram_name` varchar(255) NOT NULL,
    `memo` text,
    `start_date` datetime NOT NULL,
    `status` int NOT NULL DEFAULT '0',
    `delete_hash` tinyint DEFAULT '0',
    `customer_type` int NOT NULL DEFAULT '0',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +migrate Down
DROP TABLE `wordpress_instagrams`;