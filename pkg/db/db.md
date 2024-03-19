# Data Structures

## Long term data normalization & pinning

Extensions
```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

### NodeData

Type
```sql
CREATE TYPE nodeDataType AS (
    accumulatedUptime INTEGER,
    currentUptime INTEGER,
    firstJoined DATETIME,
    isActive BOOLEAN,
    isStaked BOOLEAN,
    isWriterNode BOOLEAN,
    lastJoined DATETIME,
    peerId VARCHAR,
);
```

Sequence

```sql
CREATE SEQUENCE nodeData_id_seq;
```

Table

```sql
CREATE TABLE nodeData
nodes (
  nodeId INTEGER PRIMARY KEY DEFAULT nextval('nodeData_id_seq'),
  node nodeDataType,
)
```

Insert (example)

```sql
INSERT INTO nodes (node) VALUES 
(ROW(
  667718450333,
  123,
  '2024-03-12T04:22:37.4547667-04:00',
  true,
  false,
  false,
  2024-03-15T09:49:57.774547667-04:00,
  '16Uiu2HAmTHk1nxbU74Co5vfbxEv56HfvPDRcJCdvYNZkGAHHQyJs'
)::nodeDataType);
```

Query (example)
```sql
SELECT (node).isStaked FROM nodes WHERE (node).peerId = '16Uiu2HAmTHk1nxbU74Co5vfbxEv56HfvPDRcJCdvYNZkGAHHQyJs'
```

### AI Sentiment Data

Sequences
```sql
CREATE SEQUENCE IF NOT EXISTS sentiment_id_seq;
CREATE SEQUENCE IF NOT EXISTS prompt_id_seq;
```

Tables
```sql
CREATE TABLE IF NOT EXISTS "public"."sentiment" (
    "id" int8 NOT NULL DEFAULT nextval('sentiment_id_seq'::regclass),
    "conversationId" int8 NOT NULL,
    "tweet" json,
    "promptId" int8
);
ALTER TABLE "public"."sentiment" ADD CONSTRAINT "sentiment_pkey" PRIMARY KEY ("id");

CREATE TABLE IF NOT EXISTS "public"."prompt" (
    "id" int8 NOT NULL DEFAULT nextval('prompt_id_seq'::regclass),
    "input" varchar(255) COLLATE "pg_catalog"."default"
);
ALTER TABLE "public"."prompt" ADD CONSTRAINT "prompt_pkey" PRIMARY KEY ("id");
```