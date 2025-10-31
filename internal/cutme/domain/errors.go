package domain

import "errors"

var (
	ErrCodeNotFound = errors.New("code not found")
	ErrCodeExpired  = errors.New("code expired")
)
