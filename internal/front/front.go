package front

import (
	"fmt"
	"net/http"
)

// Front ...
type Front struct {
	fs http.Handler
}

// New ...
func New(opt ...Option) (*Front, error) {
	opts := options{
		fs: defaultFileServer(),
	}

	for _, o := range opt {
		if err := o(&opts); err != nil {
			return nil, fmt.Errorf("front: %s", err)
		}
	}

	f := &Front{
		fs: opts.fs,
	}

	return f, nil
}

// ServeHTTP ...
func (f *Front) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.fs.ServeHTTP(w, r)
}
