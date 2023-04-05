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

	"talkee/core"
)

func newReward(db *gorm.DB, opts ...gen.DOOption) reward {
	_reward := reward{}

	_reward.rewardDo.UseDB(db, opts...)
	_reward.rewardDo.UseModel(&core.Reward{})

	tableName := _reward.rewardDo.TableName()
	_reward.ALL = field.NewAsterisk(tableName)
	_reward.ID = field.NewUint64(tableName, "id")
	_reward.TipID = field.NewUint64(tableName, "tip_id")
	_reward.ObjectType = field.NewString(tableName, "object_type")
	_reward.ObjectID = field.NewUint64(tableName, "object_id")
	_reward.SiteID = field.NewUint64(tableName, "site_id")
	_reward.RecipientID = field.NewString(tableName, "recipient_id")
	_reward.TraceID = field.NewString(tableName, "trace_id")
	_reward.SnapshotID = field.NewString(tableName, "snapshot_id")
	_reward.AssetID = field.NewString(tableName, "asset_id")
	_reward.Amount = field.NewField(tableName, "amount")
	_reward.Memo = field.NewString(tableName, "memo")
	_reward.Status = field.NewString(tableName, "status")
	_reward.CreatedAt = field.NewTime(tableName, "created_at")
	_reward.UpdatedAt = field.NewTime(tableName, "updated_at")

	_reward.fillFieldMap()

	return _reward
}

type reward struct {
	rewardDo

	ALL         field.Asterisk
	ID          field.Uint64
	TipID       field.Uint64
	ObjectType  field.String
	ObjectID    field.Uint64
	SiteID      field.Uint64
	RecipientID field.String
	TraceID     field.String
	SnapshotID  field.String
	AssetID     field.String
	Amount      field.Field
	Memo        field.String
	Status      field.String
	CreatedAt   field.Time
	UpdatedAt   field.Time

	fieldMap map[string]field.Expr
}

func (r reward) Table(newTableName string) *reward {
	r.rewardDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r reward) As(alias string) *reward {
	r.rewardDo.DO = *(r.rewardDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *reward) updateTableName(table string) *reward {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewUint64(table, "id")
	r.TipID = field.NewUint64(table, "tip_id")
	r.ObjectType = field.NewString(table, "object_type")
	r.ObjectID = field.NewUint64(table, "object_id")
	r.SiteID = field.NewUint64(table, "site_id")
	r.RecipientID = field.NewString(table, "recipient_id")
	r.TraceID = field.NewString(table, "trace_id")
	r.SnapshotID = field.NewString(table, "snapshot_id")
	r.AssetID = field.NewString(table, "asset_id")
	r.Amount = field.NewField(table, "amount")
	r.Memo = field.NewString(table, "memo")
	r.Status = field.NewString(table, "status")
	r.CreatedAt = field.NewTime(table, "created_at")
	r.UpdatedAt = field.NewTime(table, "updated_at")

	r.fillFieldMap()

	return r
}

func (r *reward) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *reward) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 14)
	r.fieldMap["id"] = r.ID
	r.fieldMap["tip_id"] = r.TipID
	r.fieldMap["object_type"] = r.ObjectType
	r.fieldMap["object_id"] = r.ObjectID
	r.fieldMap["site_id"] = r.SiteID
	r.fieldMap["recipient_id"] = r.RecipientID
	r.fieldMap["trace_id"] = r.TraceID
	r.fieldMap["snapshot_id"] = r.SnapshotID
	r.fieldMap["asset_id"] = r.AssetID
	r.fieldMap["amount"] = r.Amount
	r.fieldMap["memo"] = r.Memo
	r.fieldMap["status"] = r.Status
	r.fieldMap["created_at"] = r.CreatedAt
	r.fieldMap["updated_at"] = r.UpdatedAt
}

func (r reward) clone(db *gorm.DB) reward {
	r.rewardDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r reward) replaceDB(db *gorm.DB) reward {
	r.rewardDo.ReplaceDB(db)
	return r
}

type rewardDo struct{ gen.DO }

type IRewardDo interface {
	WithContext(ctx context.Context) IRewardDo

	SumRewardsByAsset(ctx context.Context) (result []*core.SumRewardItem, err error)
	CreateReward(ctx context.Context, model *core.Reward) (err error)
	UpdateReward(ctx context.Context, model *core.Reward) (err error)
	FindCreatedRewards(ctx context.Context, limit int) (result []*core.Reward, err error)
	GetRewardsByTipIDAndStatus(ctx context.Context, tipID uint64, status string) (result []*core.Reward, err error)
	FindRewardsByCommentIDs(ctx context.Context, commentIDs []uint64) (result []*core.Reward, err error)
}

// SELECT
//
//	"asset_id",
//	SUM("amount") as "amount"
//
// FROM
//
//	"rewards"
//
// GROUP BY "asset_id";
func (r rewardDo) SumRewardsByAsset(ctx context.Context) (result []*core.SumRewardItem, err error) {
	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT \"asset_id\", SUM(\"amount\") as \"amount\" FROM \"rewards\" GROUP BY \"asset_id\"; ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Raw(generateSQL.String()).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// INSERT INTO "rewards"
//
//		(
//	  "tip_id",
//			"object_type",
//			"object_id",
//			"site_id",
//			"recipient_id",
//			"trace_id",
//			"snapshot_id",
//			"asset_id",
//			"amount",
//	  "memo",
//			"status",
//			"created_at",
//			"updated_at"
//		)
//
// VALUES
//
//		(
//			@model.TipID,
//			@model.ObjectType,
//			@model.ObjectID,
//			@model.SiteID,
//			@model.RecipientID,
//			@model.TraceID,
//			@model.SnapshotID,
//			@model.AssetID,
//			@model.Amount,
//	  @model.Memo,
//			@model.Status,
//			NOW(), NOW()
//		);
func (r rewardDo) CreateReward(ctx context.Context, model *core.Reward) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, model.TipID)
	params = append(params, model.ObjectType)
	params = append(params, model.ObjectID)
	params = append(params, model.SiteID)
	params = append(params, model.RecipientID)
	params = append(params, model.TraceID)
	params = append(params, model.SnapshotID)
	params = append(params, model.AssetID)
	params = append(params, model.Amount)
	params = append(params, model.Memo)
	params = append(params, model.Status)
	generateSQL.WriteString("INSERT INTO \"rewards\" ( \"tip_id\", \"object_type\", \"object_id\", \"site_id\", \"recipient_id\", \"trace_id\", \"snapshot_id\", \"asset_id\", \"amount\", \"memo\", \"status\", \"created_at\", \"updated_at\" ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW() ); ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// UPDATE
//
//	"rewards"
//
// SET
//
//	"status" = @model.Status,
//	"snapshot_id"= @model.SnapshotID,
//	"trace_id"= @model.TraceID,
//	"updated_at" = NOW()
//
// WHERE
//
//	"id" = @model.ID;
func (r rewardDo) UpdateReward(ctx context.Context, model *core.Reward) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, model.Status)
	params = append(params, model.SnapshotID)
	params = append(params, model.TraceID)
	params = append(params, model.ID)
	generateSQL.WriteString("UPDATE \"rewards\" SET \"status\" = ?, \"snapshot_id\"= ?, \"trace_id\"= ?, \"updated_at\" = NOW() WHERE \"id\" = ?; ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT
// *
// FROM
// "rewards"
// WHERE
// "status" = 'created'
// ORDER BY "id" asc
// LIMIT @limit;
func (r rewardDo) FindCreatedRewards(ctx context.Context, limit int) (result []*core.Reward, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, limit)
	generateSQL.WriteString("SELECT * FROM \"rewards\" WHERE \"status\" = 'created' ORDER BY \"id\" asc LIMIT ?; ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT
// *
// FROM
// "rewards"
// WHERE
// "tip_id" = @tipID AND "status" = @status
// ORDER BY "id" asc;
// ;
func (r rewardDo) GetRewardsByTipIDAndStatus(ctx context.Context, tipID uint64, status string) (result []*core.Reward, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, tipID)
	params = append(params, status)
	generateSQL.WriteString("SELECT * FROM \"rewards\" WHERE \"tip_id\" = ? AND \"status\" = ? ORDER BY \"id\" asc; ; ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT
// *
// FROM
// "rewards"
// WHERE
// "object_id" in (@commentIDs)
// AND
// "object_type" = 'comment'
// ORDER BY "id" asc;
func (r rewardDo) FindRewardsByCommentIDs(ctx context.Context, commentIDs []uint64) (result []*core.Reward, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, commentIDs)
	generateSQL.WriteString("SELECT * FROM \"rewards\" WHERE \"object_id\" in (?) AND \"object_type\" = 'comment' ORDER BY \"id\" asc; ")

	var executeSQL *gorm.DB
	executeSQL = r.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (r rewardDo) WithContext(ctx context.Context) IRewardDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r *rewardDo) withDO(do gen.Dao) *rewardDo {
	r.DO = *do.(*gen.DO)
	return r
}
