package test

import "route256/loms/internal/model"

func compareStocks(a, b []model.Stock) bool {
	if len(a) != len(b) {
		return false
	}
	stockMap := make(map[model.SKUType]model.Stock)
	for _, stock := range a {
		stockMap[stock.SKU] = stock
	}
	for _, stock := range b {
		if s, ok := stockMap[stock.SKU]; !ok || s != stock {
			return false
		}
	}
	return true
}
