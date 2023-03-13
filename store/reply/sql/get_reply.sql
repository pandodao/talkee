SELECT
  "replies"."id",
  "replies"."user_id",
  "replies"."comment_id",
  "replies"."content",

  "users"."mixin_user_id",
  "users"."mixin_identity_number",
  "users"."full_name",
  "users"."avatar_url",

  "replies"."created_at",
  "replies"."updated_at"
FROM "replies"
INNER JOIN "users" ON "replies"."user_id" = "users"."id"
WHERE
  "replies"."id" = :reply_id
  AND "replies"."deleted_at" IS NULL
LIMIT 1
;