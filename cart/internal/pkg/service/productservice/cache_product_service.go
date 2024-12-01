package productservice

import (
	"context"
	"route256/cart/internal/logger"
	"route256/cart/internal/metrics"
	"route256/cart/internal/pkg/model"
	"sync"
)

type productServiceInterface interface {
	GetProductInfo(ctx context.Context, sku model.SKU) (*model.Product, error)
}

type ProductCacher interface {
	Put(ctx context.Context, key model.SKU, value *model.Product)
	Get(ctx context.Context, key model.SKU) (*model.Product, error)
}

type CacheProductService struct {
	origService productServiceInterface
	cache       ProductCacher
	mx          sync.Mutex
	skuMutex    map[model.SKU]*sync.Mutex
}

var _ productServiceInterface = (*CacheProductService)(nil)

func NewCacheProductService(_ int, origService productServiceInterface, casherImpl ProductCacher) *CacheProductService {
	return &CacheProductService{
		origService: origService,
		cache:       casherImpl,
		skuMutex:    map[model.SKU]*sync.Mutex{},
		mx:          sync.Mutex{},
	}
}

func (c *CacheProductService) GetProductInfo(ctx context.Context, sku model.SKU) (*model.Product, error) {
	logger.Debugw(ctx, "Getting product info from cache", "sku", sku)
	c.mx.Lock()
	mutex, ok := c.skuMutex[sku]
	if !ok {
		mutex = &sync.Mutex{}
		c.skuMutex[sku] = mutex
	}
	mutex.Lock()
	defer mutex.Unlock()
	c.mx.Unlock()
	productFromCache, err := c.cache.Get(ctx, sku)
	if err == nil {
		metrics.RecordCacheHit()
		return productFromCache, nil
	}
	metrics.RecordCacheMiss()
	logger.Debugw(ctx, "Product info not found in cache, getting from original service", "sku", sku)

	productValue, err := c.origService.GetProductInfo(ctx, sku)

	if err == nil {
		c.cache.Put(ctx, sku, productValue)
	}
	return productValue, err
}
