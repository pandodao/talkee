UPDATE
  "users"
SET
  "full_name" = :full_name,
  "avatar_url" = :avatar_url,
  "lang" = :lang,
  "updated_at" = NOW()
WHERE
  "id" = :id
;