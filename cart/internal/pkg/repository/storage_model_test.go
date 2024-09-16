package repository

import (
	"route256/cart/internal/pkg/apperrors"
	"route256/cart/internal/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	defaultItem = model.CartItem{
		UserId: 1,
		SKU:    101,
		Count:  2,
	}
)

func TestStorage_AddItem_Native(t *testing.T) {
	storage := NewStorage()

	storage.AddItem(defaultItem.UserId, defaultItem)

	assert.Equal(t, 1, len(storage.data))
	assert.Equal(t, 1, len(storage.data[defaultItem.UserId]))
	assert.Equal(t, defaultItem, storage.data[defaultItem.UserId][defaultItem.SKU])

	storage.AddItem(defaultItem.UserId, model.CartItem{UserId: 1, SKU: 101, Count: 3})

	assert.Equal(t, 1, len(storage.data))
	assert.Equal(t, 1, len(storage.data[defaultItem.UserId]))
	assert.Equal(t, uint16(5), storage.data[defaultItem.UserId][defaultItem.SKU].Count)
}

func TestStorage_RemoveItem_Native(t *testing.T) {
	storage := NewStorage()

	storage.AddItem(defaultItem.UserId, defaultItem)

	// Удаляем товар
	err := storage.RemoveItem(defaultItem.UserId, defaultItem.SKU)
	assert.NoError(t, err)

	// Проверяем внутреннее состояние
	assert.Equal(t, 0, len(storage.data[defaultItem.UserId]))
}

func TestStorage_RemoveByUserId_Native(t *testing.T) {
	storage := NewStorage()

	// Добавляем товар
	storage.AddItem(defaultItem.UserId, defaultItem)

	// Удаляем корзину пользователя
	err := storage.RemoveByUserId(defaultItem.UserId)
	assert.NoError(t, err)

	// Проверяем, что пользователь удален
	assert.Equal(t, 0, len(storage.data))
}

func TestStorage_GetCart_Native(t *testing.T) {
	storage := NewStorage()

	// Проверяем получение корзины для несуществующего пользователя
	_, err := storage.GetCart(999)
	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrCartNotFound, err)

	// Добавляем товар
	storage.AddItem(defaultItem.UserId, defaultItem)

	// Проверяем получение корзины для существующего пользователя
	cart, err := storage.GetCart(defaultItem.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cart))
	assert.Equal(t, defaultItem, cart[defaultItem.SKU])
}
