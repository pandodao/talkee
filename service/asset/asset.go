package asset

import (
	"context"
	"talkee/core"
	"talkee/store"
	"talkee/store/asset"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/patrickmn/go-cache"
)

func New(client *mixin.Client, assets core.AssetStore) *assetService {
	c := cache.New(10*time.Minute, 10*time.Minute)
	return &assetService{client: client, assets: assets, cache: c}
}

type assetService struct {
	client *mixin.Client
	assets core.AssetStore
	cache  *cache.Cache
}

func (s *assetService) UpdateAssets(ctx context.Context) error {
	as, err := s.assets.GetAssets(ctx)
	if err != nil {
		return err
	}
	var newAs []*mixin.Asset
	for _, a := range as {
		ma, err := s.client.ReadAsset(ctx, a.AssetID)
		if err != nil {
			continue
		}
		newAs = append(newAs, ma)
		s.cache.Set(ma.AssetID, *convertAsset(ma), cache.NoExpiration)
	}

	if err := store.Transaction(func(tx *store.Handler) error {
		assetStore := asset.New(tx)
		for _, as := range newAs {
			assetStore.SetAsset(ctx, convertAsset(as))
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *assetService) GetCachedAssets(ctx context.Context, assetID string) (*core.Asset, error) {
	val, found := s.cache.Get(assetID)
	if !found {
		asset, _ := s.assets.GetAsset(ctx, assetID)
		if asset == nil {
			ma, err := s.client.ReadAsset(ctx, assetID)
			if err != nil {
				return nil, err
			}
			asset = convertAsset(ma)
			s.assets.SetAsset(ctx, asset)
		}
		s.cache.Set(assetID, *asset, cache.NoExpiration)
		return asset, nil
	}
	asset := val.(core.Asset)
	return &asset, nil
}

func convertAssets(items []*mixin.Asset) []*core.Asset {
	var assets = make([]*core.Asset, len(items))
	for i, item := range items {
		assets[i] = convertAsset(item)
	}
	return assets
}

func convertAsset(item *mixin.Asset) *core.Asset {
	return &core.Asset{
		AssetID:  item.AssetID,
		Name:     item.Name,
		Symbol:   item.Symbol,
		IconURL:  item.IconURL,
		PriceUSD: item.PriceUSD,
	}
}
