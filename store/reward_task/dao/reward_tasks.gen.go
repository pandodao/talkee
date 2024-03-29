// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"talkee/core"
)

func newRewardTask(db *gorm.DB, opts ...gen.DOOption) rewardTask {
	_rewardTask := rewardTask{}

	_rewardTask.rewardTaskDo.UseDB(db, opts...)
	_rewardTask.rewardTaskDo.UseModel(&core.RewardTask{})

	tableName := _rewardTask.rewardTaskDo.TableName()
	_rewardTask.ALL = field.NewAsterisk(tableName)
	_rewardTask.ID = field.NewUint64(tableName, "id")
	_rewardTask.SiteID = field.NewUint64(tableName, "site_id")
	_rewardTask.Slug = field.NewString(tableName, "slug")
	_rewardTask.Processed = field.NewBool(tableName, "processed")
	_rewardTask.StrategyID = field.NewUint64(tableName, "strategy_id")
	_rewardTask.CreatedAt = field.NewTime(tableName, "created_at")
	_rewardTask.UpdatedAt = field.NewTime(tableName, "updated_at")

	_rewardTask.fillFieldMap()

	return _rewardTask
}

type rewardTask struct {
	rewardTaskDo

	ALL        field.Asterisk
	ID         field.Uint64
	SiteID     field.Uint64
	Slug       field.String
	Processed  field.Bool
	StrategyID field.Uint64
	CreatedAt  field.Time
	UpdatedAt  field.Time

	fieldMap map[string]field.Expr
}

func (r rewardTask) Table(newTableName string) *rewardTask {
	r.rewardTaskDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r rewardTask) As(alias string) *rewardTask {
	r.rewardTaskDo.DO = *(r.rewardTaskDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *rewardTask) updateTableName(table string) *rewardTask {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewUint64(table, "id")
	r.SiteID = field.NewUint64(table, "site_id")
	r.Slug = field.NewString(table, "slug")
	r.Processed = field.NewBool(table, "processed")
	r.StrategyID = field.NewUint64(table, "strategy_id")
	r.CreatedAt = field.NewTime(table, "created_at")
	r.UpdatedAt = field.NewTime(table, "updated_at")

	r.fillFieldMap()

	return r
}

func (r *rewardTask) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *rewardTask) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 7)
	r.fieldMap["id"] = r.ID
	r.fieldMap["site_id"] = r.SiteID
	r.fieldMap["slug"] = r.Slug
	r.fieldMap["processed"] = r.Processed
	r.fieldMap["strategy_id"] = r.StrategyID
	r.fieldMap["created_at"] = r.CreatedAt
	r.fieldMap["updated_at"] = r.UpdatedAt
}

func (r rewardTask) clone(db *gorm.DB) rewardTask {
	r.rewardTaskDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r rewardTask) replaceDB(db *gorm.DB) rewardTask {
	r.rewardTaskDo.ReplaceDB(db)
	return r
}

type rewardTaskDo struct{ gen.DO }

type IRewardTaskDo interface {
	WithContext(ctx context.Context) IRewardTaskDo
}

func (r rewardTaskDo) WithContext(ctx context.Context) IRewardTaskDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r *rewardTaskDo) withDO(do gen.Dao) *rewardTaskDo {
	r.DO = *do.(*gen.DO)
	return r
}
