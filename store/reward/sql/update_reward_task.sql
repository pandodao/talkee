UPDATE
  "reward_tasks"
SET
  "processed" = $1,
  "updated_at" = NOW()
WHERE
  "id" = $2
;