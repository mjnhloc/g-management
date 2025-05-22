CREATE TABLE `classes` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL,
    `trainer_id` BIGINT NOT NULL,
    `schedule` DATETIME NOT NULL,
    `duration` INT(5) UNSIGNED NOT NULL,
    `max_capacity` INT(3) UNSIGNED NOT NULL,
    `description` TEXT,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    CONSTRAINT `classes_fk_1` FOREIGN KEY (`trainer_id`) REFERENCES `trainers`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
