package qr

import "errors"

var (
	// ErrVersionNotFound qr version not found
	ErrVersionNotFound = errors.New("version not found")

	// ErrTooLargeSize input text is too large to encode
	ErrTooLargeSize = errors.New("data is too large to encode")
)
