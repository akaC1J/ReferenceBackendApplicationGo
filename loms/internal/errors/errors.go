package apperrors

import "errors"

var (
	ErrNotFound          = errors.New("item not found")
	ErrStockInsufficient = errors.New("insufficient stock")
	ErrNegativeReserved  = errors.New("reserved quantity cannot be negative")
	ErrNegativeAvailable = errors.New("available quantity cannot be negative")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInternal          = errors.New("internal error")
	ErrOrderState        = errors.New("invalid order state")
)
