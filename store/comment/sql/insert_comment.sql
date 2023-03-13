INSERT INTO "comments"
  (
    "user_id",
    "site_id",
    "slug",
    "content",
    "created_at", "updated_at"
  )
VALUES
  (
    :user_id,
    :site_id,
    :slug,
    :content,
    NOW(), NOW()
  )
RETURNING id
;