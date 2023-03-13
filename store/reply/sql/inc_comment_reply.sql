UPDATE
  "comments"
SET
  "reply_count" = "reply_count" + 1,
  "updated_at" = NOW()
WHERE
  "id" = :comment_id
;