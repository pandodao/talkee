BEGIN;

CREATE TABLE "favourites" (
  "id" BIGSERIAL  PRIMARY KEY,
  "comment_id" bigint,
  "user_id" bigint,

  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE UNIQUE INDEX uidx_favourite_comment_user ON "favourites" USING BTREE("comment_id", "user_id");

COMMIT;
