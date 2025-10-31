package domain

import "errors"

var (
	ErrNotFound = errors.New("link not found")
	ErrExpired  = errors.New("link expired")
)
