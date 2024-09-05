package model

type SKU int64
type UserId int64

type CartItem struct {
	SKU    SKU
	UserId UserId
	Count  uint16
}

type Product struct {
	Name  string
	Price uint32
}
