// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"

	"talkee/core"
)

func newFavourite(db *gorm.DB, opts ...gen.DOOption) favourite {
	_favourite := favourite{}

	_favourite.favouriteDo.UseDB(db, opts...)
	_favourite.favouriteDo.UseModel(&core.Favourite{})

	tableName := _favourite.favouriteDo.TableName()
	_favourite.ALL = field.NewAsterisk(tableName)
	_favourite.ID = field.NewUint64(tableName, "id")
	_favourite.UserID = field.NewUint64(tableName, "user_id")
	_favourite.CommentID = field.NewUint64(tableName, "comment_id")
	_favourite.CreatedAt = field.NewTime(tableName, "created_at")
	_favourite.UpdatedAt = field.NewTime(tableName, "updated_at")
	_favourite.DeletedAt = field.NewTime(tableName, "deleted_at")

	_favourite.fillFieldMap()

	return _favourite
}

type favourite struct {
	favouriteDo

	ALL       field.Asterisk
	ID        field.Uint64
	UserID    field.Uint64
	CommentID field.Uint64
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Time

	fieldMap map[string]field.Expr
}

func (f favourite) Table(newTableName string) *favourite {
	f.favouriteDo.UseTable(newTableName)
	return f.updateTableName(newTableName)
}

func (f favourite) As(alias string) *favourite {
	f.favouriteDo.DO = *(f.favouriteDo.As(alias).(*gen.DO))
	return f.updateTableName(alias)
}

func (f *favourite) updateTableName(table string) *favourite {
	f.ALL = field.NewAsterisk(table)
	f.ID = field.NewUint64(table, "id")
	f.UserID = field.NewUint64(table, "user_id")
	f.CommentID = field.NewUint64(table, "comment_id")
	f.CreatedAt = field.NewTime(table, "created_at")
	f.UpdatedAt = field.NewTime(table, "updated_at")
	f.DeletedAt = field.NewTime(table, "deleted_at")

	f.fillFieldMap()

	return f
}

func (f *favourite) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := f.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (f *favourite) fillFieldMap() {
	f.fieldMap = make(map[string]field.Expr, 6)
	f.fieldMap["id"] = f.ID
	f.fieldMap["user_id"] = f.UserID
	f.fieldMap["comment_id"] = f.CommentID
	f.fieldMap["created_at"] = f.CreatedAt
	f.fieldMap["updated_at"] = f.UpdatedAt
	f.fieldMap["deleted_at"] = f.DeletedAt
}

func (f favourite) clone(db *gorm.DB) favourite {
	f.favouriteDo.ReplaceConnPool(db.Statement.ConnPool)
	return f
}

func (f favourite) replaceDB(db *gorm.DB) favourite {
	f.favouriteDo.ReplaceDB(db)
	return f
}

type favouriteDo struct{ gen.DO }

type IFavouriteDo interface {
	WithContext(ctx context.Context) IFavouriteDo

	FindAllFavourites(ctx context.Context, commentID uint64) (result []*core.Favourite, err error)
	FindUserFavourites(ctx context.Context, userID uint64, commentIDs []uint64) (result []*core.Favourite, err error)
	CountAllFavourites(ctx context.Context) (result uint64, err error)
	CreateFavourite(ctx context.Context, userID uint64, commentID uint64) (err error)
	DeleteFavourite(ctx context.Context, userID uint64, commentID uint64) (err error)
}

// SELECT
//
//	*
//
// FROM "favourites"
// WHERE
//
//	"comment_id" = @commentID
//
// AND
//
//	"deleted_at" is NULL
//
// ;
func (f favouriteDo) FindAllFavourites(ctx context.Context, commentID uint64) (result []*core.Favourite, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, commentID)
	generateSQL.WriteString("SELECT * FROM \"favourites\" WHERE \"comment_id\" = ? AND \"deleted_at\" is NULL ; ")

	var executeSQL *gorm.DB
	executeSQL = f.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT
//
//	*
//
// FROM "favourites"
// WHERE
// "comment_id" IN (@commentIDs)
// AND
//
//	"user_id"=@userID
//
// AND
//
//	"deleted_at" is NULL;
func (f favouriteDo) FindUserFavourites(ctx context.Context, userID uint64, commentIDs []uint64) (result []*core.Favourite, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, commentIDs)
	params = append(params, userID)
	generateSQL.WriteString("SELECT * FROM \"favourites\" WHERE \"comment_id\" IN (?) AND \"user_id\"=? AND \"deleted_at\" is NULL; ")

	var executeSQL *gorm.DB
	executeSQL = f.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT
//
//	count("id")
//
// FROM @@table
func (f favouriteDo) CountAllFavourites(ctx context.Context) (result uint64, err error) {
	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT count(\"id\") FROM favourites ")

	var executeSQL *gorm.DB
	executeSQL = f.UnderlyingDB().Raw(generateSQL.String()).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// INSERT INTO "favourites"
//
//	(
//		"comment_id",
//		"user_id",
//		"created_at",
//		"updated_at"
//	)
//
// VALUES
//
//	(
//		@commentID,
//		@userID,
//		NOW(),
//		NOW()
//	)
//
// ON CONFLICT ("comment_id", "user_id") DO
//
//	UPDATE
//	SET "deleted_at" = NULL, "updated_at" = NOW()
//
// ;
func (f favouriteDo) CreateFavourite(ctx context.Context, userID uint64, commentID uint64) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, commentID)
	params = append(params, userID)
	generateSQL.WriteString("INSERT INTO \"favourites\" ( \"comment_id\", \"user_id\", \"created_at\", \"updated_at\" ) VALUES ( ?, ?, NOW(), NOW() ) ON CONFLICT (\"comment_id\", \"user_id\") DO UPDATE SET \"deleted_at\" = NULL, \"updated_at\" = NOW() ; ")

	var executeSQL *gorm.DB
	executeSQL = f.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// UPDATE "favourites"
// {{set}}
//
//	"deleted_at" = NOW()
//
// {{end}}
// WHERE
//
//	"user_id" = @userID
//
// AND
//
//	"comment_id" = @commentID
//
// AND
//
//	"deleted_at" is NULL
//
// ;
func (f favouriteDo) DeleteFavourite(ctx context.Context, userID uint64, commentID uint64) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("UPDATE \"favourites\" ")
	var setSQL0 strings.Builder
	setSQL0.WriteString("\"deleted_at\" = NOW() ")
	helper.JoinSetBuilder(&generateSQL, setSQL0)
	params = append(params, userID)
	params = append(params, commentID)
	generateSQL.WriteString("WHERE \"user_id\" = ? AND \"comment_id\" = ? AND \"deleted_at\" is NULL ; ")

	var executeSQL *gorm.DB
	executeSQL = f.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (f favouriteDo) WithContext(ctx context.Context) IFavouriteDo {
	return f.withDO(f.DO.WithContext(ctx))
}

func (f *favouriteDo) withDO(do gen.Dao) *favouriteDo {
	f.DO = *do.(*gen.DO)
	return f
}