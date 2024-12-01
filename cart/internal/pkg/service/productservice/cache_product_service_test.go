package productservice

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"math/rand"
	"route256/cart/internal/infra/cache/lru_cache"
	"route256/cart/internal/logger"
	"route256/cart/internal/pkg/model"
	"strconv"
	"sync"
	"testing"
	"time"
)

func (c *CacheProductService) GetProductInfoWithMx(ctx context.Context, sku model.SKU) (*model.Product, error) {
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
		return productFromCache, nil
	}
	logger.Debugw(ctx, "Product info not found in cache, getting from original service", "sku", sku)

	productValue, err := c.origService.GetProductInfo(ctx, sku)
	c.cache.Put(ctx, sku, productValue)
	return productValue, err
}

func (c *CacheProductService) GetProductInfoWithoutMx(ctx context.Context, sku model.SKU) (*model.Product, error) {
	logger.Debugw(ctx, "Getting product info from cache", "sku", sku)
	productFromCache, err := c.cache.Get(ctx, sku)
	if err == nil {
		return productFromCache, nil
	}
	logger.Debugw(ctx, "Product info not found in cache, getting from original service", "sku", sku)

	productValue, err := c.origService.GetProductInfo(ctx, sku)
	c.cache.Put(ctx, sku, productValue)
	return productValue, err
}

// Бенчмарк с большим JSON
func BenchmarkCacheProductService_GetProductInfoWithMx(b *testing.B) {
	logger.NewNopLogger()
	mc := minimock.NewController(b)
	origService := NewProductServiceInterfaceMock(mc)
	origService.GetProductInfoMock.Set(func(ctx context.Context, sku model.SKU) (*model.Product, error) {
		// Имитация задержки оригинального сервиса
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return &model.Product{
			Name:  strconv.Itoa(int(sku)),
			Price: uint32(sku * 3),
		}, nil
	})
	ctx := context.Background()

	// Параметры для многопоточной загрузки
	const numThreads = 20           // Количество потоков (горутины)
	const requestsPerThread = 10000 // Количество запросов на поток

	b.Run("WithMutex", func(b *testing.B) {
		cacheProductService := NewCacheProductService(0, origService, lru_cache.NewLruCache[model.SKU, *model.Product](100000))

		b.ResetTimer() // Сброс таймера перед началом теста
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for t := 0; t < numThreads; t++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < requestsPerThread; j++ {
						sku := model.SKU(j % 10)
						_, _ = cacheProductService.GetProductInfoWithMx(ctx, sku)
					}
				}()
			}
			wg.Wait() // Ждём завершения всех потоков
		}
	})

	b.Run("WithoutMutex", func(b *testing.B) {
		cacheProductService := NewCacheProductService(0, origService, lru_cache.NewLruCache[model.SKU, *model.Product](10000))

		b.ResetTimer() // Сброс таймера перед началом теста
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for t := 0; t < numThreads; t++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < requestsPerThread; j++ {
						sku := model.SKU(j % 10)
						_, _ = cacheProductService.GetProductInfoWithoutMx(ctx, sku)
					}
				}()
			}
			wg.Wait() // Ждём завершения всех потоков
		}
	})
}
