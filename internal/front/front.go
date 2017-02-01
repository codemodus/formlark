package front

import (
	"fmt"
	"net/http"
)

// Front ...
type Front struct{}

// New ...
func New() (*Front, error) {
	f := &Front{}

	return f, nil
}

// ServeHTTP ...
func (f *Front) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "front")
}
