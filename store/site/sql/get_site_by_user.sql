SELECT
  "id",
  "user_id",
  "origins",
  "name",
  "use_arweave",
  "reward_strategy",
  "created_at",
  "updated_at"
FROM "sites"
WHERE
  "user_id" = :user_id
  AND
  deleted_at IS NULL
;