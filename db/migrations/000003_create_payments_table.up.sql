CREATE TABLE `payments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `member_id` BIGINT NOT NULL,
    `membership_id` BIGINT NOT NULL,
    `amount` INT(11) UNSIGNED NOT NULL,
    `payment_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `payment_method` ENUM('cash', 'credit_card', 'bank_transfer') DEFAULT 'cash',
    `status` ENUM('completed', 'failed', 'refundeded') DEFAULT 'completed',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    CONSTRAINT `payments_fk_1` FOREIGN KEY (`member_id`) REFERENCES `members`(`id`),
    CONSTRAINT `payments_fk_2` FOREIGN KEY (`membership_id`) REFERENCES `memberships`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
