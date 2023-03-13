INSERT INTO assets ("asset_id", "name", "symbol", "icon_url", "price_usd", "created_at", "updated_at")
VALUES (:asset_id, :name, :symbol, :icon_url, :price_usd, NOW(), NOW())
ON CONFLICT (asset_id) DO
  UPDATE SET
    price_usd=EXCLUDED.price_usd,
    name=EXCLUDED.name,
    symbol=EXCLUDED.symbol,
    icon_url=EXCLUDED.icon_url,
    updated_at=NOW();
