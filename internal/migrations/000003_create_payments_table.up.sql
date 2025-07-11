CREATE TABLE `payments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `membership_id` BIGINT NOT NULL,
    `price` INT(11) UNSIGNED NOT NULL,
    `payment_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `payment_method` ENUM('cash', 'credit_card', 'bank_transfer') DEFAULT 'cash',
    `status` ENUM('completed', 'failed', 'refunded') DEFAULT 'completed',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    CONSTRAINT `payments_fk_1` FOREIGN KEY (`membership_id`) REFERENCES `memberships`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
