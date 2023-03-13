UPDATE
  "favourites"
SET
  "deleted_at" = NOW()
WHERE
  "comment_id" = :comment_id
  AND "user_id" = :user_id
;