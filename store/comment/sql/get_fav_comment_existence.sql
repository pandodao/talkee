-- check the existence
SELECT
  id, comment_id, user_id
FROM "favourites"
WHERE
  "comment_id" = :comment_id
  AND "user_id" = :user_id
  AND "deleted_at" IS NULL
LIMIT 1;