UPDATE
  "rewards"
SET
  "status" = :status,
  "snapshot_id"=:snapshot_id,
  "trace_id"=:trace_id,
  "updated_at" = NOW()
WHERE
  "id" = :id
;