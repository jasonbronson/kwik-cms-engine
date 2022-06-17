
DROP TABLE IF EXISTS "role";
CREATE TABLE "role" ("id" VARCHAR NOT NULL UNIQUE,"name" VARCHAR UNIQUE,"created_at" timestamp,"updated_at" timestamp,PRIMARY KEY ("id"));

ALTER TABLE "public"."user" 
  ADD "role_id" varchar, 
  ADD "first_name" VARCHAR,
  ADD "last_name" VARCHAR, 
  ADD "title" VARCHAR,
  ADD "content" VARCHAR, 
  ADD "facebook_url" VARCHAR,
  ADD "instagram_url" VARCHAR,
  ADD "twitter_url" VARCHAR,
  ADD "youtube_url" VARCHAR;