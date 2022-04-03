CREATE TABLE IF NOT EXISTS `tokens` (
    `id` INT AUTO_INCREMENT NOT NULL,
    `user_id` INT NOT NULL,
    `uuid` VARCHAR (255) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_tokens_on_user_id` (`user_id`),
    INDEX `index_tokens_on_uuid` (`uuid`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `authorities` (
    `id` INT AUTO_INCREMENT NOT NULL,
    `token_id` INT NOT NULL,
    `type` VARCHAR (255) NOT NULL,
    `right` BIT (64) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `index_authorities_on_token_id` (`token_id`),
    INDEX `index_authorities_on_type` (`type`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;
