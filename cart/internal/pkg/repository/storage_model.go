package repository

import (
	"route256/cart/internal/pkg/apperrors"
	"route256/cart/internal/pkg/model"
)

type Storage struct {
	data map[model.UserId]map[model.SKU]model.CartItem
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[model.UserId]map[model.SKU]model.CartItem),
	}
}

func (s *Storage) AddItem(userID model.UserId, item model.CartItem) {
	if _, ok := s.data[userID]; !ok {
		s.data[userID] = make(map[model.SKU]model.CartItem)
	}
	if existingItem, ok := s.data[userID][item.SKU]; ok {
		existingItem.Count += item.Count
		s.data[userID][item.SKU] = existingItem
	} else {
		s.data[userID][item.SKU] = item
	}
}

func (s *Storage) RemoveItem(userID model.UserId, sku model.SKU) error {
	if items, ok := s.data[userID]; ok {
		delete(items, sku)
		if len(items) == 0 {
			delete(s.data, userID) // Удаляем пользователя, если его корзина пуста
		}
		return nil
	}
	return apperrors.ErrCartNotFound
}

func (s *Storage) RemoveByUserId(userID model.UserId) error {
	if _, ok := s.data[userID]; ok {
		delete(s.data, userID)
		return nil
	}
	return apperrors.ErrUserNotFound
}

func (s *Storage) GetCart(userID model.UserId) (map[model.SKU]model.CartItem, error) {
	cart, ok := s.data[userID]
	if !ok {
		return nil, apperrors.ErrCartNotFound
	}
	return cart, nil
}
