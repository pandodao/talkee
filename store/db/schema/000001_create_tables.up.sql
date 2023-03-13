BEGIN;

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

CREATE INDEX idx_site_user ON "sites" USING BTREE("user_id");
CREATE UNIQUE INDEX uidx_site_origin ON "sites" USING BTREE("origins");

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
CREATE INDEX idx_comment_general_created ON "comments" USING BTREE("site_id", "slug", "created_at");
CREATE INDEX idx_comment_general_favor ON "comments" USING BTREE("site_id", "slug", "favor_count");

CREATE TABLE "replies" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint,
  "comment_id" bigint,
  "content" varchar(5120) NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE INDEX idx_reply_comment ON "replies" USING BTREE("comment_id");

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

  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE INDEX idx_users_mixin_user_id ON "users" USING BTREE("mixin_user_id");

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

COMMIT;
