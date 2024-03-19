CREATE TABLE IF NOT EXISTS "public"."sentiment" (
    "id" int8 NOT NULL DEFAULT nextval('sentiment_id_seq'::regclass),
    "conversation_id" int8 NOT NULL,
    "tweet" json,
    "prompt_id" int8
);
