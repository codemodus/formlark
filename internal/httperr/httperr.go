package httperr

import (
	"fmt"
	"net/http"
)

// HTTPError ...
type HTTPError interface {
	error
	Err() error
	Status() int
}

// HTTPErr ...
type HTTPErr struct {
	error
	code      int
	sanitized string
}

// New ...
func New(err error, code int, sanitized string) *HTTPErr {
	e := &HTTPErr{
		error:     err,
		code:      code,
		sanitized: sanitized,
	}

	return e
}

// Err ...
func (e *HTTPErr) Err() error {
	if e.error == nil {
		return fmt.Errorf(e.Error())
	}

	return e.error
}

// Status ...
func (e *HTTPErr) Status() int {
	if e.code < 100 || e.code > 600 {
		return http.StatusInternalServerError
	}

	return e.code
}

// Error ...
func (e *HTTPErr) Error() string {
	if e.sanitized == "" {
		return http.StatusText(e.Status())
	}

	return e.sanitized
}
