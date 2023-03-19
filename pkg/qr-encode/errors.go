package qr_encode

import "errors"

var (
	ErrVersionNotFound = errors.New("version not found")

	ErrTooLargeSize = errors.New("data is too large to encode")
)
