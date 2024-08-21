package errorz

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidRole     = errors.New("invalid role")
)
