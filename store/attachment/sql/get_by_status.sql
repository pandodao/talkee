SELECT
  *
FROM attachments
WHERE
  "status" = :status
  AND deleted_at IS NULL
ORDER BY
  created_at
LIMIT :limit
;