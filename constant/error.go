package constant

import "errors"

var (
	ErrMethodNotAllowed    = errors.New("method not allowed")
	ErrInvalidContentType  = errors.New("invalid content type")
	ErrEmptyBody           = errors.New("request body is empty")
	ErrInvalidJSON         = errors.New("invalid JSON")
	ErrMultipleJSONObject  = errors.New("request body must contain only one JSON object")
	ErrRequestBodyTooLarge = errors.New("request body too large")
)
