package front

import (
	"fmt"
	"net/http"
)

// Front ...
type Front struct{}

// ServeHTTP ...
func (f *Front) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "front")
}
