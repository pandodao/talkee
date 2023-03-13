SELECT
  "key",
  "value",
  "updated_at"
FROM
  "properties"
WHERE
  "key" = :key
LIMIT
  1;