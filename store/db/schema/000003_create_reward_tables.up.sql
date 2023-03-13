BEGIN;

CREATE TABLE "rewards" (
  "id" BIGSERIAL PRIMARY KEY,
  "object_type" varchar(256),
  "object_id" bigint  NOT NULL,
  "site_id" int  NOT NULL,
  "recipient_id" varchar(36) NOT NULL,
  "trace_id" varchar(36),
  "snapshot_id" varchar(36),
  "asset_id" varchar(36) NOT NULL,
  "amount" numeric(64,8) NOT NULL,
  "status" varchar(50) NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz
);


CREATE TABLE "reward_tasks" (
  "id" BIGSERIAL PRIMARY KEY,
  "site_id" int  NOT NULL,
  "slug" varchar(1000) NOT NULL,
  "processed" boolean default false,
  "strategy_id" int default 0,
  "created_at" timestamptz,
  "updated_at" timestamptz
);


CREATE UNIQUE INDEX uidx_reward_tasks_site_slug ON "reward_tasks" USING BTREE("site_id", "slug");


CREATE TABLE "reward_strategies" (
  "id" BIGSERIAL PRIMARY KEY,
  "params" text  NOT NULL,
  "is_default" boolean default false
);




COMMIT;