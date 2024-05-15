CREATE TABLE "public"."work" (
    "id" int8 NOT NULL DEFAULT nextval('work_id_seq'::regclass),
    "uuid" uuid NOT NULL,
    "payload" jsonb NOT NULL,
    "response" jsonb
);