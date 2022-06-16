DROP TABLE IF EXISTS "author";
CREATE TABLE `author` (`id` VARCHAR NOT NULL UNIQUE,`username` VARCHAR UNIQUE,`email` VARCHAR,`first_name` VARCHAR,`last_name` VARCHAR,`title` VARCHAR,`content` VARCHAR,`facebook_url` VARCHAR,`instagram_url` VARCHAR,`twitter_url` VARCHAR,`youtube_url` VARCHAR,`slug` VARCHAR UNIQUE,`media_id` VARCHAR,`created_at` datetime,`updated_at` datetime,PRIMARY KEY (`id`),CONSTRAINT `fk_post_author` FOREIGN KEY (`id`) REFERENCES `post`(`author_id`));

DROP TABLE IF EXISTS "categories_post_links";
CREATE TABLE `categories_post_links` (`category_id` VARCHAR,`post_id` VARCHAR);

DROP TABLE IF EXISTS "category";
CREATE TABLE `category` (`id` VARCHAR NOT NULL UNIQUE,`name` VARCHAR UNIQUE,`created_at` datetime,`updated_at` datetime,PRIMARY KEY (`id`));

DROP TABLE IF EXISTS "media";
CREATE TABLE `media` (`id` VARCHAR NOT NULL UNIQUE,`name` VARCHAR,`url` VARCHAR,`alt_text` VARCHAR,`created_at` datetime,`updated_at` datetime,PRIMARY KEY (`id`),CONSTRAINT `fk_author_media` FOREIGN KEY (`id`) REFERENCES `author`(`media_id`),CONSTRAINT `fk_post_media` FOREIGN KEY (`id`) REFERENCES `post`(`media_id`));

DROP TABLE IF EXISTS "post";
CREATE TABLE `post` (`id` VARCHAR NOT NULL UNIQUE,`title` VARCHAR,`slug` VARCHAR UNIQUE,`content` TEXT,`summary` VARCHAR,`description` VARCHAR,`author_id` VARCHAR,`media_id` VARCHAR,`publish_date` TIMESTAMP,`created_at` datetime,`updated_at` datetime,PRIMARY KEY (`id`));

DROP TABLE IF EXISTS "tag";
CREATE TABLE `tag` (`id` VARCHAR NOT NULL UNIQUE,`name` VARCHAR UNIQUE,`created_at` datetime,`updated_at` datetime,PRIMARY KEY (`id`));

DROP TABLE IF EXISTS "tags_post_links";
CREATE TABLE `tags_post_links` (`tag_id` VARCHAR,`post_id` VARCHAR);

DROP TABLE IF EXISTS "user";
CREATE TABLE `user` (`id` VARCHAR NOT NULL UNIQUE,`email` VARCHAR UNIQUE,`password` VARCHAR,`created_at` datetime,`updated_at` datetime,PRIMARY KEY (`id`));

