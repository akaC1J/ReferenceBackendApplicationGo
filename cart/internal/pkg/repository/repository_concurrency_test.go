package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"route256/cart/internal/pkg/model"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentInsert(t *testing.T) {
	t.Parallel()
	storage := NewStorage()
	repo := NewRepository(storage)

	userId := model.UserId(1)
	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Одновременное добавление элементов
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			item := model.CartItem{
				UserId: userId,
				SKU:    model.SKU(i),
				Count:  1,
			}
			_, _ = repo.InsertItem(context.Background(), item)
		}(i)
	}

	wg.Wait()

	// Проверяем, что все элементы добавлены
	cart, err := repo.GetCartByUserId(context.Background(), userId)
	assert.NoError(t, err)
	assert.Equal(t, numGoroutines, len(cart))
	for i := 0; i < numGoroutines; i++ {
		sku := model.SKU(i)
		item, exists := cart[sku]
		assert.True(t, exists, "SKU %d should exist in cart", sku)
		assert.Equal(t, uint16(1), item.Count)
	}
}

func TestConcurrentGetCart(t *testing.T) {
	t.Parallel()
	storage := NewStorage()
	repo := NewRepository(storage)

	userId := model.UserId(2)
	for i := 0; i < 50; i++ {
		item := model.CartItem{
			UserId: userId,
			SKU:    model.SKU(i),
			Count:  uint16(i + 1),
		}
		_, err := repo.InsertItem(context.Background(), item)
		assert.NoError(t, err)
	}

	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Одновременное чтение
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			cart, _ := repo.GetCartByUserId(context.Background(), userId)
			assert.Equal(t, true, cartHasUserId(cart, userId))
			for sku := range cart {
				assert.Equal(t, cart[sku].UserId, userId)
			}
			assert.Equal(t, 50, len(cart))
		}()
	}

	wg.Wait()
}

func TestConcurrentInsertAndGet(t *testing.T) {
	t.Parallel()
	storage := NewStorage()
	repo := NewRepository(storage)

	userId := model.UserId(3)
	numInsertGoroutines := 50
	numGetGoroutines := 50
	var wg sync.WaitGroup
	wg.Add(numInsertGoroutines + numGetGoroutines)

	// Запускаем вставки
	for i := 0; i < numInsertGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			item := model.CartItem{
				UserId: userId,
				SKU:    model.SKU(i),
				Count:  uint16(i + 1),
			}
			_, _ = repo.InsertItem(context.Background(), item)
		}(i)
	}

	// Одновременное чтение
	for i := 0; i < numGetGoroutines; i++ {
		go func() {
			defer wg.Done()
			cart, _ := repo.GetCartByUserId(context.Background(), userId)
			assert.Equal(t, true, cartHasUserId(cart, userId))
			// Количество элементов может меняться, но не должно превышать вставленных
			assert.LessOrEqual(t, len(cart), numInsertGoroutines)
		}()
	}

	wg.Wait()

	// Финальная проверка
	cart, _ := repo.GetCartByUserId(context.Background(), userId)
	assert.Equal(t, true, cartHasUserId(cart, userId))
	assert.Equal(t, numInsertGoroutines, len(cart))
	for i := 0; i < numInsertGoroutines; i++ {
		sku := model.SKU(i)
		item, exists := cart[sku]
		assert.True(t, exists, "SKU %d should exist in cart", sku)
		assert.Equal(t, uint16(i+1), item.Count)
	}
}

func TestConcurrentInsertAndRemove(t *testing.T) {
	t.Parallel()
	storage := NewStorage()
	repo := NewRepository(storage)

	userId := model.UserId(4)
	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	// Запускаем вставки
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Microsecond)
			item := model.CartItem{
				UserId: userId,
				SKU:    model.SKU(i),
				Count:  1,
			}
			_, err := repo.InsertItem(context.Background(), item)
			assert.NoError(t, err)
		}(i)
	}

	removeSuccessCount := atomic.Int32{}

	// Запускаем удаления
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			err := repo.RemoveItem(context.Background(), userId, model.SKU(i))
			if err == nil {
				removeSuccessCount.Add(1)
			}
		}(i)
	}

	wg.Wait()
	cart, _ := repo.GetCartByUserId(context.Background(), userId)
	//количество ошибок + количество найденных пользователей должно быть равно количеству
	//вставленных элементов или количеству горутин
	assert.Equal(t, numGoroutines, len(cart)+int(removeSuccessCount.Load()))
}

func cartHasUserId(cart map[model.SKU]model.CartItem, userId model.UserId) bool {
	for sku := range cart {
		if cart[sku].UserId != userId {
			return false
		}
	}
	return true
}
