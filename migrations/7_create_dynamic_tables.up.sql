CREATE TABLE "public"."dynamic_groups_fields" (
    "dynamic_field_id" varchar,
    "dynamic_group_id" varchar
);

CREATE SEQUENCE IF NOT EXISTS dynamic_groups_id_seq;

-- Table Definition
CREATE TABLE "public"."dynamic_groups" (
    "id" varchar NOT NULL DEFAULT nextval('dynamic_groups_id_seq'::regclass),
    "title" varchar,
    "created_at" timestamp,
    "updated_at" timestamp,
    PRIMARY KEY ("id")
);

CREATE SEQUENCE IF NOT EXISTS dynamic_fields_id_seq;

-- Table Definition
CREATE TABLE "public"."dynamic_fields" (
    "id" varchar NOT NULL DEFAULT nextval('dynamic_fields_id_seq'::regclass),
    "label" varchar,
    "name" varchar,
    "type" varchar,
    "instructions" varchar,
    "created_at" timestamp,
    "updated_at" timestamp,
    PRIMARY KEY ("id")
);