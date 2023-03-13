SELECT
  "snapshot_id",
  "trace_id",
  "source",
  "transaction_hash",
  "receiver",
  "sender",
  "type",
  "created_at",
  "user_id",
  "opponent_id",
  "asset_id",
  "amount",
  "memo"
FROM snapshots
WHERE
  snapshot_id = :id
LIMIT 1
;