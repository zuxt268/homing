
-- +migrate Up
CREATE TABLE IF NOT EXISTS `posts` (
    `id` int NOT NULL AUTO_INCREMENT,
    `media_id` varchar(45) NOT NULL,
    `customer_id` int NOT NULL,
    `timestamp` varchar(45) DEFAULT NULL,
    `media_url` mediumtext,
    `permalink` varchar(255) DEFAULT NULL,
    `wordpress_link` varchar(255) DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=139 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +migrate Down
DROP TABLE `posts`;