
-- +migrate Up
CREATE TABLE IF NOT EXISTS `wordpress_instagrams` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `wordpress` varchar(255) NOT NULL,
    `instagram_id` varchar(255) NOT NULL,
    `memo` text,
    `start_date` datetime NOT NULL,
    `status` int NOT NULL DEFAULT '0',
    `delete_hash` tinyint DEFAULT '0',
    `customer_type` int NOT NULL DEFAULT '0',
    `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +migrate Down
DROP TABLE `wordpress_instagrams`;