INSERT INTO "properties"
  ("key", "value", "updated_at", "deleted_at")
VALUES
  (:key, :value, NOW(), NULL)
ON CONFLICT ("key") DO
  UPDATE SET
    "value"=EXCLUDED.value,
    updated_at=NOW();
;