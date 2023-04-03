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

func newUser(db *gorm.DB, opts ...gen.DOOption) user {
	_user := user{}

	_user.userDo.UseDB(db, opts...)
	_user.userDo.UseModel(&core.User{})

	tableName := _user.userDo.TableName()
	_user.ALL = field.NewAsterisk(tableName)
	_user.ID = field.NewUint64(tableName, "id")
	_user.MixinUserID = field.NewString(tableName, "mixin_user_id")
	_user.MixinIdentityNumber = field.NewString(tableName, "mixin_identity_number")
	_user.FullName = field.NewString(tableName, "full_name")
	_user.AvatarURL = field.NewString(tableName, "avatar_url")
	_user.MvmPublicKey = field.NewString(tableName, "mvm_public_key")
	_user.Lang = field.NewString(tableName, "lang")
	_user.CreatedAt = field.NewTime(tableName, "created_at")
	_user.UpdatedAt = field.NewTime(tableName, "updated_at")
	_user.DeletedAt = field.NewTime(tableName, "deleted_at")

	_user.fillFieldMap()

	return _user
}

type user struct {
	userDo

	ALL                 field.Asterisk
	ID                  field.Uint64
	MixinUserID         field.String
	MixinIdentityNumber field.String
	FullName            field.String
	AvatarURL           field.String
	MvmPublicKey        field.String
	Lang                field.String
	CreatedAt           field.Time
	UpdatedAt           field.Time
	DeletedAt           field.Time

	fieldMap map[string]field.Expr
}

func (u user) Table(newTableName string) *user {
	u.userDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u user) As(alias string) *user {
	u.userDo.DO = *(u.userDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *user) updateTableName(table string) *user {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewUint64(table, "id")
	u.MixinUserID = field.NewString(table, "mixin_user_id")
	u.MixinIdentityNumber = field.NewString(table, "mixin_identity_number")
	u.FullName = field.NewString(table, "full_name")
	u.AvatarURL = field.NewString(table, "avatar_url")
	u.MvmPublicKey = field.NewString(table, "mvm_public_key")
	u.Lang = field.NewString(table, "lang")
	u.CreatedAt = field.NewTime(table, "created_at")
	u.UpdatedAt = field.NewTime(table, "updated_at")
	u.DeletedAt = field.NewTime(table, "deleted_at")

	u.fillFieldMap()

	return u
}

func (u *user) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *user) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 10)
	u.fieldMap["id"] = u.ID
	u.fieldMap["mixin_user_id"] = u.MixinUserID
	u.fieldMap["mixin_identity_number"] = u.MixinIdentityNumber
	u.fieldMap["full_name"] = u.FullName
	u.fieldMap["avatar_url"] = u.AvatarURL
	u.fieldMap["mvm_public_key"] = u.MvmPublicKey
	u.fieldMap["lang"] = u.Lang
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["updated_at"] = u.UpdatedAt
	u.fieldMap["deleted_at"] = u.DeletedAt
}

func (u user) clone(db *gorm.DB) user {
	u.userDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u user) replaceDB(db *gorm.DB) user {
	u.userDo.ReplaceDB(db)
	return u
}

type userDo struct{ gen.DO }

type IUserDo interface {
	WithContext(ctx context.Context) IUserDo

	GetUser(ctx context.Context, id uint64) (result *core.User, err error)
	GetUserByIDs(ctx context.Context, ids []uint64) (result []*core.User, err error)
	GetUserByMixinID(ctx context.Context, mixinUserID string) (result *core.User, err error)
	CreateUser(ctx context.Context, user *core.User) (result uint64, err error)
	UpdateBasicInfo(ctx context.Context, id uint64, user *core.User) (err error)
}

// SELECT
//
//	*
//
// FROM @@table
// WHERE id = @id AND deleted_at IS NULL;
func (u userDo) GetUser(ctx context.Context, id uint64) (result *core.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM users WHERE id = ? AND deleted_at IS NULL; ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT
//
//	*
//
// FROM @@table
// WHERE id in (@ids) AND deleted_at IS NULL;
func (u userDo) GetUserByIDs(ctx context.Context, ids []uint64) (result []*core.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, ids)
	generateSQL.WriteString("SELECT * FROM users WHERE id in (?) AND deleted_at IS NULL; ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT
//
//	*
//
// FROM @@table
// WHERE mixin_user_id = @mixinUserID;
func (u userDo) GetUserByMixinID(ctx context.Context, mixinUserID string) (result *core.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, mixinUserID)
	generateSQL.WriteString("SELECT * FROM users WHERE mixin_user_id = ?; ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// INSERT INTO users
//
//	(
//		"full_name", "avatar_url",
//		"mixin_user_id", "mixin_identity_number",
//		"lang", "mvm_public_key",
//		"created_at", "updated_at"
//	)
//
// VALUES
//
//	(
//		@user.FullName, @user.AvatarURL,
//		@user.MixinUserID, @user.MixinIdentityNumber,
//		@user.Lang, @user.MvmPublicKey,
//		NOW(), NOW()
//	)
//
// RETURNING id;
func (u userDo) CreateUser(ctx context.Context, user *core.User) (result uint64, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, user.FullName)
	params = append(params, user.AvatarURL)
	params = append(params, user.MixinUserID)
	params = append(params, user.MixinIdentityNumber)
	params = append(params, user.Lang)
	params = append(params, user.MvmPublicKey)
	generateSQL.WriteString("INSERT INTO users ( \"full_name\", \"avatar_url\", \"mixin_user_id\", \"mixin_identity_number\", \"lang\", \"mvm_public_key\", \"created_at\", \"updated_at\" ) VALUES ( ?, ?, ?, ?, ?, ?, NOW(), NOW() ) RETURNING id; ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// UPDATE @@table
//
//	{{set}}
//	  "name"=@user.FullName,
//	  "avatar_url"=@user.AvatarURL,
//	  "lang"=@user.Lang,
//		"updated_at"=NOW()
//	{{end}}
//
// WHERE
//
//	"id" = @id AND "deleted_at" IS NULL;
func (u userDo) UpdateBasicInfo(ctx context.Context, id uint64, user *core.User) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("UPDATE users ")
	var setSQL0 strings.Builder
	params = append(params, user.FullName)
	params = append(params, user.AvatarURL)
	params = append(params, user.Lang)
	setSQL0.WriteString("\"name\"=?, \"avatar_url\"=?, \"lang\"=?, \"updated_at\"=NOW() ")
	helper.JoinSetBuilder(&generateSQL, setSQL0)
	params = append(params, id)
	generateSQL.WriteString("WHERE \"id\" = ? AND \"deleted_at\" IS NULL; ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (u userDo) WithContext(ctx context.Context) IUserDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u *userDo) withDO(do gen.Dao) *userDo {
	u.DO = *do.(*gen.DO)
	return u
}