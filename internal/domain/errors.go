package domain

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidOrder = errors.New("invalid order")
)
