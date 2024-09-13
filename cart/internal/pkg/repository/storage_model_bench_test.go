package repository

import (
	"route256/cart/internal/pkg/apperrors"
	"testing"

	"route256/cart/internal/pkg/model"
)

func BenchmarkStorage_AddItem(b *testing.B) {
	s := NewStorage()

	item := model.CartItem{
		SKU:   model.SKU(1),
		Count: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uid := model.UserId(i % 1000)
		sku := model.SKU(i % 100)
		item.SKU = sku

		s.AddItem(uid, item)
	}
}

func BenchmarkStorage_RemoveItem(b *testing.B) {
	s := NewStorage()

	userID := model.UserId(1)
	item := model.CartItem{
		SKU:   model.SKU(1),
		Count: 1,
	}

	for i := 0; i < b.N; i++ {
		sku := model.SKU(i % 100)
		item.SKU = sku
		s.AddItem(userID, item)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sku := model.SKU(i % 100)
		err := s.RemoveItem(userID, sku)
		if err != nil && err != apperrors.ErrCartNotFound {
			b.Errorf("failed to remove item: %v", err)
		}
	}
}
