
-- +migrate Up
CREATE TABLE IF NOT EXISTS `google_posts` (
    `id` int NOT NULL AUTO_INCREMENT,
    `google_business_url` varchar(500) NOT NULL,
    `instagram_url` varchar(500) NOT NULL,
    `media_id` varchar(255) NOT NULL,
    `customer_id` int NOT NULL,
    `name` varchar(500) NOT NULL,
    `media_format` varchar(50) NOT NULL,
    `google_url` varchar(500) NOT NULL,
    `create_time` varchar(255) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +migrate Down
DROP TABLE `google_posts`;