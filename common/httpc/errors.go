package httpc

import "errors"

var (
	ErrInvalidBaseURL       = errors.New("invalid base url")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)
