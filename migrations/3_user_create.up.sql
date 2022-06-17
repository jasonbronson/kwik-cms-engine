CREATE SEQUENCE IF NOT EXISTS media_id_seq;

CREATE TABLE "public"."media" (
    "id" int4 NOT NULL DEFAULT nextval('media_id_seq'::regclass),
    "name" varchar,
    "size" numeric,
    "url" varchar,
    "created_at" timestamp,
    "updated_at" timestamp,
    "content" text,
    "title" text,
    "excerpt" text,
    PRIMARY KEY ("id")
);

ALTER TABLE "public"."post" RENAME COLUMN "author_id" TO "user_id";
