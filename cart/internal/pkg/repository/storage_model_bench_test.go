package repository

import (
	"testing"

	"route256/cart/internal/pkg/model"
)

func BenchmarkStorage_AddItem(b *testing.B) {
	s := NewStorage()

	items := make([]model.CartItem, 100)
	uids := make([]model.UserId, 1000)

	for i := range items {
		items[i] = model.CartItem{
			SKU:   model.SKU(i),
			Count: uint16(i % 10),
		}
	}

	for i := range uids {
		uids[i] = model.UserId(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.AddItem(uids[i%1000], items[i%100])
	}
}

func BenchmarkStorage_RemoveItem(b *testing.B) {
	s := NewStorage()

	userID := model.UserId(1)

	items := make([]model.CartItem, 100)
	skus := make([]model.SKU, 100)

	for i := range items {
		skus[i] = model.SKU(i)
		items[i] = model.CartItem{
			SKU:   skus[i],
			Count: 1,
		}
	}

	for i := 0; i < b.N; i++ {
		s.AddItem(userID, items[i%100])
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.RemoveItem(userID, skus[i%100])
	}
}
