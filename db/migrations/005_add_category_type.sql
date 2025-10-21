-- +migrate Up
ALTER TABLE `Category` ADD COLUMN `category_type` int NOT NULL DEFAULT 0;

-- category_type = 1: 収入
UPDATE `Category` SET `category_type` = 1 WHERE `id` IN (1, 2, 3);

-- category_type = 2: 支出
UPDATE `Category` SET `category_type` = 2 WHERE `id` IN (5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19);

-- category_type = 3: 貯金
UPDATE `Category` SET `category_type` = 3 WHERE `id` IN (20, 21);

-- category_type = 4: NISA
UPDATE `Category` SET `category_type` = 4 WHERE `id` IN (22, 23);

-- +migrate Down
ALTER TABLE `Category` DROP COLUMN `category_type`;
