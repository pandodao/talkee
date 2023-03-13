SELECT
  COUNT("comments"."id")
FROM "comments"
WHERE
  "comments"."deleted_at" IS NULL
;