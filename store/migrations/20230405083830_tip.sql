-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "tips" (
  "id" BIGSERIAL PRIMARY KEY,
  "uuid" varchar(36) NOT NULL,
  "user_id" bigint NOT NULL,
  "site_id" bigint NOT NULL,
  "slug" varchar(255) NOT NULL,
  "tip_type" varchar(255) NOT NULL,
  "strategy_name" varchar(255) NOT NULL,
  "strategy_params" jsonb DEFAULT '{}'::jsonb,
  "asset_id" varchar(36) NOT NULL,
  "amount" numeric(64,8),
  "memo" varchar(255) DEFAULT '',
  "status" int default 0,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);
CREATE INDEX "idx_tip_uuid" ON "tips" ("uuid", "deleted_at");
CREATE INDEX "idx_tip_status" ON "tips" ("status", "deleted_at");

ALTER TABLE "rewards" ADD COLUMN "tip_id" bigint NOT NULL;
ALTER TABLE "rewards" ADD COLUMN "memo" varchar(255) NOT NULL;
CREATE INDEX "idx_reward_tip_id_status" ON "rewards" ("tip_id", "status");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "tips";
ALTER TABLE "rewards" DROP COLUMN IF EXISTS "tip_id";
ALTER TABLE "rewards" DROP COLUMN IF EXISTS "memo";
DROP INDEX IF EXISTS "idx_reward_tip_id_status";
-- +goose StatementEnd
