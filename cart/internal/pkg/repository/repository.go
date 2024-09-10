package repository

import (
	"context"
	"route256/cart/internal/pkg/model"
)

type AbstractStorage interface {
	AddItem(id model.UserId, item model.CartItem)
	RemoveItem(id model.UserId, sku model.SKU) error
	RemoveByUserId(id model.UserId) error
	GetCart(id model.UserId) (map[model.SKU]model.CartItem, error)
}

// Repository использует Storage для работы с корзинами
type Repository struct {
	storage AbstractStorage
}

// NewRepository создает новый репозиторий с хранилищем корзин
func NewRepository(storage AbstractStorage) *Repository {
	return &Repository{storage: storage}
}

// InsertItem добавляет или обновляет элемент в корзине пользователя
func (r *Repository) InsertItem(_ context.Context, cartItem model.CartItem) (*model.CartItem, error) {
	r.storage.AddItem(cartItem.UserId, cartItem)
	return &cartItem, nil
}

// RemoveItem удаляет товар из корзины пользователя
func (r *Repository) RemoveItem(_ context.Context, userId model.UserId, sku model.SKU) error {
	return r.storage.RemoveItem(userId, sku)
}

// RemoveByUserId удаляет корзину пользователя
func (r *Repository) RemoveByUserId(_ context.Context, userId model.UserId) error {
	return r.storage.RemoveByUserId(userId)
}

// GetItem возвращает корзину пользователя
func (r *Repository) GetItem(_ context.Context, userId model.UserId) (map[model.SKU]model.CartItem, error) {
	return r.storage.GetCart(userId)
}
