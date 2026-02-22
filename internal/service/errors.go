package service

import "errors"

var (
	ErrNotFound    = errors.New("product not found")
	ErrValidation  = errors.New("validation error")
)
