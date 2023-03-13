INSERT INTO attachments
(
  "size",
  "mime_type",
  "persisted",
  "duration",
  "path",
  "category",
  "provider",
  "provider_item_id",
  "status",
  "lang",
  "text",

  "created_at",
  "updated_at"
)
VALUES
(
  0,
  :mime_type,
  false,
  0,
  '',
  :category,
  :provider,
  '',
  :status,
  :lang,
  :text,
  NOW(),
  NOW()
)
RETURNING id
;