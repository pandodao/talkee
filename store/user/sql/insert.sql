INSERT INTO users
  (
    "full_name", "avatar_url",
    "mixin_user_id", "mixin_identity_number",
    "lang", "mvm_public_key",
    "created_at", "updated_at"
  )
VALUES
  (
    :full_name, :avatar_url,
    :mixin_user_id, :mixin_identity_number,
    :lang, :mvm_public_key,
    NOW(), NOW()
  )
RETURNING id
;