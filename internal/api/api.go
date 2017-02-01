package api

import (
	"fmt"
	"net/http"
)

// API ...
type API struct{}

// New ...
func New() (*API, error) {
	a := &API{}

	return a, nil
}

// ServeHTTP ...
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "api")
}
