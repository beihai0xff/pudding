CREATE TABLE IF NOT EXISTS `webhook_trigger_template` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime NULL,
    `updated_at` datetime NULL,
    `deleted_at` datetime NULL,
    `topic` varchar(255) NOT NULL DEFAULT 'unknown' COMMENT 'message topic',
    `payload` TEXT NOT NULL COMMENT 'message content',
    `deliver_after` int NOT NULL COMMENT 'message deliver after seconds',
    `excepted_end_time` TIMESTAMP NOT NULL COMMENT 'excepted trigger end time, if it is 0, it means that it will not end',
    `excepted_loop_times` int unsigned NOT NULL DEFAULT 0 COMMENT 'except loop times',
    `looped_times` int unsigned NOT NULL DEFAULT 0 COMMENT 'already loop times',
    `status` int unsigned NOT NULL DEFAULT 0 COMMENT 'trigger template status: enable->1 disable->2 offline->3 and so on',
    PRIMARY KEY (`id`),
    INDEX `idx_webhook_trigger_template_deleted_at` (`deleted_at`)
) ENGINE = InnoDB;
