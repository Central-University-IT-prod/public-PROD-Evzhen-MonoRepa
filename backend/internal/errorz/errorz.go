package errorz

import "errors"

var (
	InternalServerError = errors.New("internal server error")
	CannotParseJSON     = errors.New("cannot parse JSON")
)
