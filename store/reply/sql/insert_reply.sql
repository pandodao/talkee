INSERT INTO "replies"
  (
    "user_id",
    "comment_id",
    "content",
    "created_at", "updated_at"
  )
VALUES
  (
    :user_id,
    :comment_id,
    :content,
    NOW(), NOW()
  )
RETURNING id
;