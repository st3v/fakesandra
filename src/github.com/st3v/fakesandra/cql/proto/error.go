package proto

import "errors"

var (
	// io
	errMaxLenExceeded = errors.New("Exceeds maximum length")

	// versioning
	errUnsupportedProtocolVersion = errors.New("Unsupported CQL Protocol version")

	// routing & handling
	errMissingRoute   = errors.New("Missing route")
	errMissingHandler = errors.New("Missing handler")
)
