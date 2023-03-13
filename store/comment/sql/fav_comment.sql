INSERT INTO "favourites"
  (
    "comment_id",
    "user_id",
    "created_at",
    "updated_at"
  )
VALUES
  (
    :comment_id,
    :user_id,
    NOW(),
    NOW()
  )
ON CONFLICT ("comment_id", "user_id") DO
  UPDATE
  SET "deleted_at" = NULL, "updated_at" = NOW()
;


