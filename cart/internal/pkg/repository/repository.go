package repository

import (
	"context"
	"errors"
	"route256/cart/internal/pkg/model"
)

// storage если существует userId, то map[model.SKU]model.CartItem быть не может nil или пустой
// поэтому если userId существуе
type storage map[model.UserId]map[model.SKU]model.CartItem

// Repository Для поддержания контракта все мои функции будут возращать error(в любой момент мы сменим с in-memory на persistence
// хранилище в котором могут быть проблемы)
type Repository struct {
	storage storage
}

func NewRepository() *Repository {
	return &Repository{storage: make(storage)}
}

func (r *Repository) InsertItem(_ context.Context, cartItem model.CartItem) (*model.CartItem, error) {
	_, ok := r.storage[cartItem.UserId]
	if !ok {
		r.storage[cartItem.UserId] = make(map[model.SKU]model.CartItem)
		r.storage[cartItem.UserId][cartItem.SKU] = cartItem
		return &cartItem, nil
	}

	if item, ok := r.storage[cartItem.UserId][cartItem.SKU]; ok {
		// If the item already exists, increase the count
		item.Count += cartItem.Count
		r.storage[cartItem.UserId][cartItem.SKU] = item
		return &item, nil
	}

	// If the item does not exist, add a new one
	r.storage[cartItem.UserId][cartItem.SKU] = cartItem
	return &cartItem, nil
}

// RemoveItem поддерживаем условие, что если корзина пустая удаляем пользователя
func (r *Repository) RemoveItem(_ context.Context, userId model.UserId, sku model.SKU) error {
	if items, ok := r.storage[userId]; ok {
		delete(items, sku)
		if len(items) == 0 {
			delete(r.storage, userId)
		}
	}
	return nil
}

func (r *Repository) RemoveByUserId(_ context.Context, userId model.UserId) error {
	if _, ok := r.storage[userId]; ok {
		delete(r.storage, userId)
	}
	return nil
}

func (r *Repository) GetItem(_ context.Context, userId model.UserId) (map[model.SKU]model.CartItem, error) {
	cartUser, ok := r.storage[userId]
	if !ok {
		return nil, errors.New("cart's user not exist")
	}
	return cartUser, nil
}
