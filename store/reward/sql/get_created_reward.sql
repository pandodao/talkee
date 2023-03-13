SELECT 
"id",
"object_type",
"object_id",
"site_id",
"recipient_id",
"trace_id",
"snapshot_id",
"asset_id",
"amount",
"status",
"created_at",
"updated_at"
FROM 
"rewards"
WHERE 
"status" = 'created'

ORDER BY "id" asc

LIMIT $1
;