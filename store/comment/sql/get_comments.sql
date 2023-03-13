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
  "comments"."site_id" = :site_id
  AND "comments"."slug" = :slug
  AND "comments"."deleted_at" IS NULL
ORDER BY "%s" %s
OFFSET :offset
LIMIT :limit
;