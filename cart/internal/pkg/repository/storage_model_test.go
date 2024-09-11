package repository

import (
	"route256/cart/internal/pkg/apperrors"
	"route256/cart/internal/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_AddItem_Native(t *testing.T) {
	storage := NewStorage()

	item := model.CartItem{
		UserId: 1,
		SKU:    101,
		Count:  2,
	}

	// Добавляем новый товар
	storage.AddItem(item.UserId, item)

	// Проверяем внутреннее состояние структуры
	assert.Equal(t, 1, len(storage.data))
	assert.Equal(t, 1, len(storage.data[item.UserId]))
	assert.Equal(t, item, storage.data[item.UserId][item.SKU])

	// Добавляем тот же товар снова
	storage.AddItem(item.UserId, model.CartItem{UserId: 1, SKU: 101, Count: 3})

	// Проверяем, что количество товара увеличилось
	assert.Equal(t, 1, len(storage.data))
	assert.Equal(t, 1, len(storage.data[item.UserId]))
	assert.Equal(t, uint16(5), storage.data[item.UserId][item.SKU].Count)
}

func TestStorage_RemoveItem_Native(t *testing.T) {
	storage := NewStorage()

	item := model.CartItem{
		UserId: 1,
		SKU:    101,
		Count:  2,
	}

	// Добавляем товар
	storage.AddItem(item.UserId, item)

	// Удаляем товар
	err := storage.RemoveItem(item.UserId, item.SKU)
	assert.NoError(t, err)

	// Проверяем внутреннее состояние
	assert.Equal(t, 0, len(storage.data[item.UserId]))
}

func TestStorage_RemoveByUserId_Native(t *testing.T) {
	storage := NewStorage()

	item := model.CartItem{
		UserId: 1,
		SKU:    101,
		Count:  2,
	}

	// Добавляем товар
	storage.AddItem(item.UserId, item)

	// Удаляем корзину пользователя
	err := storage.RemoveByUserId(item.UserId)
	assert.NoError(t, err)

	// Проверяем, что пользователь удален
	assert.Equal(t, 0, len(storage.data))
}

func TestStorage_GetCart_Native(t *testing.T) {
	storage := NewStorage()

	item := model.CartItem{
		UserId: 1,
		SKU:    101,
		Count:  2,
	}

	// Проверяем получение корзины для несуществующего пользователя
	_, err := storage.GetCart(999)
	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrCartNotFound, err)

	// Добавляем товар
	storage.AddItem(item.UserId, item)

	// Проверяем получение корзины для существующего пользователя
	cart, err := storage.GetCart(item.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cart))
	assert.Equal(t, item, cart[item.SKU])
}
