SELECT
	"id",
	"comment_id",
	"user_id",
	"created_at",
	"updated_at"
FROM "favourites"
WHERE 

    "comment_id" = $1
	AND
	"deleted_at"  is NULL
;