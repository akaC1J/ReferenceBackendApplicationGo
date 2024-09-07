package apperrors

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrCartNotFound = errors.New("cart not found")
)
