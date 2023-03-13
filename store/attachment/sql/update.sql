UPDATE
  "attachments"
SET
  size = :size,
  persisted = :persisted,
  duration = :duration,
  path = :path,
  status = :status,
  provider_item_id = :provider_item_id,
  updated_at = NOW()
WHERE
  "id" = :id
;