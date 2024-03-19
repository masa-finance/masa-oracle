CREATE TABLE IF NOT EXISTS "public"."prompt" (
    "id" int8 NOT NULL DEFAULT nextval('prompt_id_seq'::regclass),
    "input" varchar(255) COLLATE "pg_catalog"."default"
);
