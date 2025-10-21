-- +migrate Up
ALTER TABLE `Category` ADD COLUMN `category_type` int NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE `Category` DROP COLUMN `category_type`;
