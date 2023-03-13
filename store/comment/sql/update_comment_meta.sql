UPDATE
  "comments"
SET
  :key = :value,
  "updated_at" = NOW()
WHERE
  "id" = :id
;