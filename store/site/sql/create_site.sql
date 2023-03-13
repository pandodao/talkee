INSERT INTO "sites"
  (
    "user_id",
    "name",
    "origins",
    "use_arweave",
    "reward_strategy",
    "created_at", "updated_at"
  )
VALUES
  (
    :user_id,
    :name,
    :origins,
    :use_arweave,
    :reward_strategy,
    NOW(), NOW()
  )
RETURNING id
;