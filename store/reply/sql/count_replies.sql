SELECT
  COUNT("id")
FROM "replies"
WHERE
  "replies"."comment_id" = :comment_id
  AND "replies"."deleted_at" IS NULL
;