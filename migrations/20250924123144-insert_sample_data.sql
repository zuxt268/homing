-- +migrate Up
-- Insert sample customers
INSERT INTO `customers` (`id`, `name`, `email`, `password`, `wordpress_url`, `facebook_token`, `start_date`, `instagram_business_account_id`, `instagram_business_account_name`, `instagram_token_status`, `delete_hash`, `payment_type`, `type`) VALUES
(1, '田中太郎', 'tanaka@example.com', '$2a$10$9iC2kGdq8XrJ5qP7dF9R.ewFEkD.zqEPcQ4jqc9xZfKrHgL8bMnCy', 'https://tanaka-blog.com', 'facebook_token_123', '2024-01-01 09:00:00', 'ig_business_123', 'tanaka_shop', 1, 0, 'monthly', 0),
(2, '山田花子', 'yamada@example.com', '$2a$10$9iC2kGdq8XrJ5qP7dF9R.ewFEkD.zqEPcQ4jqc9xZfKrHgL8bMnCy', 'https://yamada-store.com', 'facebook_token_456', '2024-02-15 10:30:00', 'ig_business_456', 'yamada_fashion', 1, 0, 'yearly', 0),
(3, '佐藤次郎', 'sato@example.com', '$2a$10$9iC2kGdq8XrJ5qP7dF9R.ewFEkD.zqEPcQ4jqc9xZfKrHgL8bMnCy', 'https://sato-cafe.com', NULL, '2024-03-10 14:20:00', 'ig_business_789', 'sato_cafe', 0, 0, 'none', 1),
(4, '鈴木美香', 'suzuki@example.com', '$2a$10$9iC2kGdq8XrJ5qP7dF9R.ewFEkD.zqEPcQ4jqc9xZfKrHgL8bMnCy', 'https://suzuki-beauty.com', 'facebook_token_789', '2024-04-05 11:15:00', 'ig_business_012', 'suzuki_beauty', 1, 0, 'monthly', 0);

-- Insert sample posts
INSERT INTO `posts` (`id`, `media_id`, `customer_id`, `timestamp`, `media_url`, `permalink`, `wordpress_link`, `created_at`) VALUES
(1, 'media_001', 1, '1640995200', 'https://instagram.com/p/sample1.jpg', 'https://instagram.com/p/sample1/', 'https://tanaka-blog.com/post/instagram-1/', '2024-01-01 12:00:00'),
(2, 'media_002', 1, '1641081600', 'https://instagram.com/p/sample2.jpg', 'https://instagram.com/p/sample2/', 'https://tanaka-blog.com/post/instagram-2/', '2024-01-02 12:00:00'),
(3, 'media_003', 2, '1644753600', 'https://instagram.com/p/sample3.jpg', 'https://instagram.com/p/sample3/', 'https://yamada-store.com/post/instagram-1/', '2024-02-15 15:30:00'),
(4, 'media_004', 2, '1644840000', 'https://instagram.com/p/sample4.jpg', 'https://instagram.com/p/sample4/', 'https://yamada-store.com/post/instagram-2/', '2024-02-16 15:30:00'),
(5, 'media_005', 3, '1647334800', 'https://instagram.com/p/sample5.jpg', 'https://instagram.com/p/sample5/', 'https://sato-cafe.com/post/instagram-1/', '2024-03-10 18:20:00'),
(6, 'media_006', 4, '1649154900', 'https://instagram.com/p/sample6.jpg', 'https://instagram.com/p/sample6/', 'https://suzuki-beauty.com/post/instagram-1/', '2024-04-05 14:15:00');

-- +migrate Down
-- Delete sample data
DELETE FROM `posts` WHERE `media_id` IN ('media_001', 'media_002', 'media_003', 'media_004', 'media_005', 'media_006');
DELETE FROM `customers` WHERE `email` IN ('tanaka@example.com', 'yamada@example.com', 'sato@example.com', 'suzuki@example.com');