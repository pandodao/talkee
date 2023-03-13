SELECT
  "asset_id",
  sum("amount") as "amount"
FROM
  "rewards"
GROUP BY "asset_id";