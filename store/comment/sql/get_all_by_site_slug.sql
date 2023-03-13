SELECT
  "id",
  "user_id",
  "site_id",
  "slug",
  "favor_count"
FROM "comments"
WHERE
"site_id" = $1 
AND
"slug" = $2
;