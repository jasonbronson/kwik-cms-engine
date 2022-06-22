DROP SEQUENCE IF EXISTS media_id_seq;
DROP TABLE IF EXISTS "public"."media";

CREATE TABLE "public"."media" (
    "id" VARCHAR NOT NULL UNIQUE,
    "name" varchar,
    "size" numeric,
    "url" varchar,
    "created_at" timestamp,
    "updated_at" timestamp,
    PRIMARY KEY ("id")
);