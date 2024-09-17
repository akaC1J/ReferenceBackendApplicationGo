package model

import "fmt"

type SKU int64
type UserId int64

type CartItem struct {
	SKU    SKU
	UserId UserId
	Count  uint16
}

type StockInfo struct {
	SKU           SKU
	TotalCount    uint32
	ReservedCount uint32
}

type Product struct {
	Name  string
	Price uint32
}

// Validate проверяет, что все поля CartItem корректны
func (ci *CartItem) Validate() error {
	if ci.SKU < 1 {
		return fmt.Errorf("SKU must be positive")
	}
	if ci.UserId < 1 {
		return fmt.Errorf("UserId must be positive")
	}
	if ci.Count < 1 {
		return fmt.Errorf("Сount must be positive")
	}
	return nil
}
