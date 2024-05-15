CREATE TABLE "public"."event" (
    "id" int8 NOT NULL DEFAULT nextval('work_id_seq'::regclass),
    "work_id" uuid NOT NULL,
    "payload" jsonb NOT NULL
);