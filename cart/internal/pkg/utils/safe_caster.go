package utils

import (
	"fmt"
	"math"
	"route256/cart/internal/pkg/model"
)

func SafeInt64ToUint32(value model.SKU) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, fmt.Errorf("value %d is out of uint32 range", value)
	}
	return uint32(value), nil
}
