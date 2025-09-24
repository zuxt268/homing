
-- +migrate Up
CREATE TABLE IF NOT EXISTS `customers` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `wordpress_url` varchar(255) NOT NULL,
    `facebook_token` varchar(255) DEFAULT NULL,
    `start_date` datetime DEFAULT NULL,
    `instagram_business_account_id` varchar(255) DEFAULT NULL,
    `instagram_business_account_name` varchar(255) DEFAULT NULL,
    `instagram_token_status` int NOT NULL DEFAULT '0',
    `delete_hash` tinyint DEFAULT '0',
    `payment_type` varchar(45) NOT NULL DEFAULT 'none',
    `type` int NOT NULL DEFAULT '0' COMMENT '顧客種別: 0=普通, 1=アウトソーシング',
    PRIMARY KEY (`id`),
    UNIQUE KEY `email_UNIQUE` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- +migrate Down
DROP TABLE `customers`;