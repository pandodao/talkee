SELECT
  "comments"."id",
  "comments"."user_id",
  "comments"."site_id",
  "comments"."slug",
  "comments"."favor_count", "reply_count",
  "comments"."arweave_tx_hash",
  "comments"."content",

  "users"."mixin_user_id",
  "users"."mixin_identity_number",
  "users"."full_name",
  "users"."avatar_url",
  "users"."mvm_public_key",

  "comments"."created_at",
  "comments"."updated_at"
FROM "comments"
INNER JOIN "users" ON "comments"."user_id" = "users"."id"
WHERE
  "comments"."site_id" in (
    select "id" FROM "sites" WHERE  "use_arweave" = true
  )
  AND  ("comments"."arweave_tx_hash" is NULL OR "comments"."arweave_tx_hash" = '')
  AND "comments"."deleted_at" IS NULL
ORDER BY "comments"."created_at" asc
LIMIT :limit
;