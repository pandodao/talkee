-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "sites" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint,
  "name" varchar(256) NOT NULL,
  "origins" text[] NOT NULL default '{}',
  "use_arweave" boolean default false,
  "reward_strategy" int default 0,

  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE INDEX idx_site_user ON "sites" USING BTREE("user_id", "deleted_at");
CREATE INDEX uidx_site_origin ON "sites" USING BTREE("origins", "deleted_at");

CREATE TABLE "comments" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint,
  "site_id" bigint,
  "slug" varchar(256),
  "favor_count" bigint default 0,
  "reply_count" bigint default 0,
  "arweave_tx_hash" varchar(128) default '',
  "content" varchar(10240) NOT NULL,

  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE INDEX idx_comment_general ON "comments" USING BTREE("site_id", "slug", "deleted_at");
CREATE INDEX idx_comment_general_created ON "comments" USING BTREE("site_id", "slug", "created_at", "deleted_at");
CREATE INDEX idx_comment_general_favor ON "comments" USING BTREE("site_id", "slug", "favor_count", "deleted_at");

CREATE TABLE "replies" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint,
  "comment_id" bigint,
  "content" varchar(5120) NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE INDEX idx_reply_comment ON "replies" USING BTREE("comment_id", "deleted_at");

CREATE TABLE "properties" (
  "key" varchar(255) PRIMARY KEY,
  "value" varchar(255),
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "mixin_user_id" uuid,
  "mixin_identity_number" varchar(16),
  "avatar_url" varchar(1024) default '',
  "full_name" varchar(128) NOT NULL,
  "lang" varchar(16) default 'en',
  "mvm_public_key" varchar(128) default '',

  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);
CREATE INDEX idx_users_mixin_user_id ON "users" USING BTREE("mixin_user_id");
CREATE INDEX idx_users_mvm_public_key ON "users" USING BTREE("mvm_public_key");

CREATE TABLE "attachments" (
  "id" BIGSERIAL PRIMARY KEY,
  "size" bigint,
  "mime_type" varchar(64) default '',
  "persisted" boolean,
  "path" text default '',
  "category" int,
  "bucket" varchar(64),
  "status" int default 0,

  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);


CREATE TABLE "assets" (
  "asset_id" char(36) PRIMARY KEY,
  "name" varchar(255),
  "symbol" varchar(255),
  "icon_url" varchar(255),
  "price_usd" numeric(64,8),
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz,
  "order" bigint DEFAULT 0
);

CREATE INDEX idx_assets_order ON "assets" USING BTREE("order");
CREATE INDEX idx_assets_symbol ON "assets" USING BTREE("symbol");


CREATE TABLE "snapshots" (
  "snapshot_id" varchar(36) PRIMARY KEY,
  "trace_id" varchar(36),
  "source" varchar(32),
  "transaction_hash" varchar(64),
  "receiver" varchar(256),
  "sender" varchar(256),
  "type" varchar(32),
  "user_id" varchar(36),
  "opponent_id" varchar(36),
  "asset_id" varchar(36),
  "memo" varchar(256),
  "amount" numeric(64,8),
  "created_at" timestamptz
);

CREATE TABLE "favourites" (
  "id" BIGSERIAL PRIMARY KEY,
  "comment_id" bigint,
  "user_id" bigint,

  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);
CREATE UNIQUE INDEX uidx_favourite_comment_user ON "favourites" USING BTREE("comment_id", "user_id", "deleted_at");


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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS properties;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS replies;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS favourites;
DROP TABLE IF EXISTS rewards;
-- +goose StatementEnd
