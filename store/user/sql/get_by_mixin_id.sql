SELECT
  "id",
  "mixin_user_id",
  "mixin_identity_number",
  "full_name",
  "avatar_url",
  "mvm_public_key",
  "lang",

  "created_at",
  "updated_at"
FROM users
WHERE mixin_user_id = :mixin_user_id;