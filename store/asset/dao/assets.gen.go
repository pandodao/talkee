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

func newAsset(db *gorm.DB, opts ...gen.DOOption) asset {
	_asset := asset{}

	_asset.assetDo.UseDB(db, opts...)
	_asset.assetDo.UseModel(&core.Asset{})

	tableName := _asset.assetDo.TableName()
	_asset.ALL = field.NewAsterisk(tableName)
	_asset.AssetID = field.NewString(tableName, "asset_id")
	_asset.Name = field.NewString(tableName, "name")
	_asset.PriceUSD = field.NewField(tableName, "price_usd")
	_asset.Symbol = field.NewString(tableName, "symbol")
	_asset.IconURL = field.NewString(tableName, "icon_url")
	_asset.Order_ = field.NewInt64(tableName, "order")
	_asset.CreatedAt = field.NewTime(tableName, "created_at")
	_asset.UpdatedAt = field.NewTime(tableName, "updated_at")
	_asset.DeletedAt = field.NewTime(tableName, "deleted_at")

	_asset.fillFieldMap()

	return _asset
}

type asset struct {
	assetDo

	ALL       field.Asterisk
	AssetID   field.String
	Name      field.String
	PriceUSD  field.Field
	Symbol    field.String
	IconURL   field.String
	Order_    field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Time

	fieldMap map[string]field.Expr
}

func (a asset) Table(newTableName string) *asset {
	a.assetDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a asset) As(alias string) *asset {
	a.assetDo.DO = *(a.assetDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *asset) updateTableName(table string) *asset {
	a.ALL = field.NewAsterisk(table)
	a.AssetID = field.NewString(table, "asset_id")
	a.Name = field.NewString(table, "name")
	a.PriceUSD = field.NewField(table, "price_usd")
	a.Symbol = field.NewString(table, "symbol")
	a.IconURL = field.NewString(table, "icon_url")
	a.Order_ = field.NewInt64(table, "order")
	a.CreatedAt = field.NewTime(table, "created_at")
	a.UpdatedAt = field.NewTime(table, "updated_at")
	a.DeletedAt = field.NewTime(table, "deleted_at")

	a.fillFieldMap()

	return a
}

func (a *asset) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *asset) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 9)
	a.fieldMap["asset_id"] = a.AssetID
	a.fieldMap["name"] = a.Name
	a.fieldMap["price_usd"] = a.PriceUSD
	a.fieldMap["symbol"] = a.Symbol
	a.fieldMap["icon_url"] = a.IconURL
	a.fieldMap["order"] = a.Order_
	a.fieldMap["created_at"] = a.CreatedAt
	a.fieldMap["updated_at"] = a.UpdatedAt
	a.fieldMap["deleted_at"] = a.DeletedAt
}

func (a asset) clone(db *gorm.DB) asset {
	a.assetDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a asset) replaceDB(db *gorm.DB) asset {
	a.assetDo.ReplaceDB(db)
	return a
}

type assetDo struct{ gen.DO }

type IAssetDo interface {
	WithContext(ctx context.Context) IAssetDo

	GetAssets(ctx context.Context) (result []*core.Asset, err error)
	GetAsset(ctx context.Context, assetID string) (result *core.Asset, err error)
	SetAsset(ctx context.Context, asset *core.Asset) (err error)
}

// SELECT
// * FROM @@table
// WHERE deleted_at IS NULL;
func (a assetDo) GetAssets(ctx context.Context) (result []*core.Asset, err error) {
	var generateSQL strings.Builder
	generateSQL.WriteString("SELECT * FROM assets WHERE deleted_at IS NULL; ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String()).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// SELECT
// * FROM @@table
// WEHRE
//
//	asset_id = @assetID AND deleted_at IS NULL;
func (a assetDo) GetAsset(ctx context.Context, assetID string) (result *core.Asset, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, assetID)
	generateSQL.WriteString("SELECT * FROM assets WEHRE asset_id = ? AND deleted_at IS NULL; ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// INSERT INTO assets ("asset_id", "name", "symbol", "icon_url", "price_usd", "created_at", "updated_at")
// VALUES (@asset.AssetID, @asset.Name, @asset.Symbol, @asset.IconURL, @asset.PriceUSD, NOW(), NOW())
// ON CONFLICT (asset_id) DO
//
//	UPDATE SET
//		price_usd=EXCLUDED.price_usd,
//		name=EXCLUDED.name,
//		symbol=EXCLUDED.symbol,
//		icon_url=EXCLUDED.icon_url,
//		updated_at=NOW();
func (a assetDo) SetAsset(ctx context.Context, asset *core.Asset) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, asset.AssetID)
	params = append(params, asset.Name)
	params = append(params, asset.Symbol)
	params = append(params, asset.IconURL)
	params = append(params, asset.PriceUSD)
	generateSQL.WriteString("INSERT INTO assets (\"asset_id\", \"name\", \"symbol\", \"icon_url\", \"price_usd\", \"created_at\", \"updated_at\") VALUES (?, ?, ?, ?, ?, NOW(), NOW()) ON CONFLICT (asset_id) DO UPDATE SET price_usd=EXCLUDED.price_usd, name=EXCLUDED.name, symbol=EXCLUDED.symbol, icon_url=EXCLUDED.icon_url, updated_at=NOW(); ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (a assetDo) WithContext(ctx context.Context) IAssetDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a *assetDo) withDO(do gen.Dao) *assetDo {
	a.DO = *do.(*gen.DO)
	return a
}
