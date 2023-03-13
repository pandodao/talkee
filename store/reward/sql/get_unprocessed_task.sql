-- check the existence
SELECT
"id",
"site_id",
"slug",
"processed",
"created_at",
"updated_at"
FROM "reward_tasks"
WHERE
  "created_at" < $1
  AND "processed" = false
LIMIT $2;