SELECT
  "id"
FROM "reward_tasks"
WHERE
  site_id = $1 AND slug = $2
LIMIT 1;