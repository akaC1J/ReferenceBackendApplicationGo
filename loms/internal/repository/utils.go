package repository

import (
	"fmt"
	"math"
)

func SafeInt64ToUint32(value int64) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, fmt.Errorf("value %d is out of uint32 range", value)
	}
	return uint32(value), nil
}
