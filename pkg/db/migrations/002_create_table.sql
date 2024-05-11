CREATE TABLE "public"."work" (
    "id" int8 NOT NULL DEFAULT nextval('work_id_seq'::regclass),
    "payload" jsonb NOT NULL,
    "raw" jsonb
);