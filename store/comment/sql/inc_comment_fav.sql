UPDATE
  "comments"
SET
  "favor_count" = "favor_count" + 1,
  "updated_at" = NOW()
WHERE
  "id" = :comment_id
;