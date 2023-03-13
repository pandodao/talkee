SELECT
  *
FROM attachments
WHERE
  "provider" = :provider
  "provider_item_id" = :provider_item_id
;
