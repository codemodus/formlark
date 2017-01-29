package api

import (
	"fmt"
	"net/http"
)

// API ...
type API struct{}

// ServeHTTP ...
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "api")
}
