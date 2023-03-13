SELECT
	"id",
	"comment_id"
FROM "favourites"
WHERE 
    "comment_id" in(?)
	AND
	"user_id"  = ?
	AND
	"deleted_at"  is NULL
;