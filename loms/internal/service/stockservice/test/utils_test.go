package test

import "route256/loms/internal/model"

func compareStocks(a, b map[model.SKUType]*model.Stock) bool {
	if len(a) != len(b) {
		return false
	}
	for sku, stock := range a {
		stock2, ok := b[sku]
		if !ok {
			return false
		}
		if stock.SKU != stock2.SKU {
			return false
		}
		if stock.TotalCount != stock2.TotalCount {
			return false
		}
		if stock.ReservedCount != stock2.ReservedCount {
			return false
		}
	}
	return true
}
