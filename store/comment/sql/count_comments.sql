SELECT
  COUNT("comments"."id")
FROM "comments"
WHERE
  "comments"."site_id" = :site_id
  AND "comments"."slug" = :slug
  AND "comments"."deleted_at" IS NULL
;