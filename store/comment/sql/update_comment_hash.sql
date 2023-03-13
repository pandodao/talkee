UPDATE
  "comments"
SET
  "arweave_tx_hash" = :arweave_tx_hash,
  "updated_at" = NOW()
WHERE
  "id" = :id
;