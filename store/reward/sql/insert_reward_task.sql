INSERT INTO "reward_tasks"
  (
"site_id", 
"slug",
"processed",
"strategy_id",
"created_at",
"updated_at"
  )
VALUES
  (
    $1,
    $2,
    false,
    $3,
    NOW(), NOW()
  )
RETURNING id
;