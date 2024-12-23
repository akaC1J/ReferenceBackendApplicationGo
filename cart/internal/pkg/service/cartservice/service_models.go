package cartservice

import "route256/cart/internal/pkg/model"

type EnrichedCartItem struct {
	SKU   int64  `json:"sku_id"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type CartContent struct {
	Items      []EnrichedCartItem `json:"items"`
	TotalPrice uint32             `json:"total_price"`
}

func createEnrichedCartItemDTO(cartItem model.CartItem, product model.Product) EnrichedCartItem {
	return EnrichedCartItem{
		SKU:   int64(cartItem.SKU),
		Count: cartItem.Count,
		Name:  product.Name,
		Price: product.Price,
	}
}
