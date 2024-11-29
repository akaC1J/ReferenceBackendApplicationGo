package productservice

import (
	"context"
	"route256/cart/internal/infra/cache"
	"route256/cart/internal/logger"
	"route256/cart/internal/pkg/model"
)

type productServiceInterface interface {
	GetProductInfo(ctx context.Context, sku model.SKU) (*model.Product, error)
}

type CacheProductService struct {
	origService productServiceInterface
	cache       *cache.LruCache[model.SKU, *model.Product]
}

var _ productServiceInterface = (*CacheProductService)(nil)

func NewCacheProductService(capacity int, origService productServiceInterface) *CacheProductService {
	return &CacheProductService{
		origService: origService,
		cache:       cache.NewLruCache[model.SKU, *model.Product](capacity),
	}
}

func (c *CacheProductService) GetProductInfo(ctx context.Context, sku model.SKU) (*model.Product, error) {
	logger.Debugw(ctx, "Getting product info from cache", "sku", sku)
	productFromCache, ok := c.cache.Get(sku)
	if ok {
		return productFromCache, nil
	}
	logger.Debugw(ctx, "Product info not found in cache, getting from original service", "sku", sku)
	productValue, err := c.origService.GetProductInfo(ctx, sku)
	c.cache.Put(sku, productValue)
	return productValue, err
}
