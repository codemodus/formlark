package dommux

import (
	"fmt"
	"net/http"
)

// DomMux ...
type DomMux struct {
	m map[string]http.Handler
	h http.Handler
}

// New ...
func New(opt ...Option) (*DomMux, error) {
	opts := options{
		m: make(map[string]http.Handler),
	}

	for _, o := range opt {
		if err := o(&opts); err != nil {
			return nil, fmt.Errorf("dommux: %s", err)
		}
	}

	if len(opts.m) == 0 {
		return nil, fmt.Errorf("dommux: must provide at least one domain/handler")
	}

	d := &DomMux{
		m: make(map[string]http.Handler),
		h: http.NotFoundHandler(),
	}

	for k, v := range opts.m {
		d.m[k] = v
	}

	if opts.h != nil {
		d.h = opts.h
	}

	return d, nil
}

// ServeHTTP ...
func (d *DomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, ok := d.m[r.Host]
	if !ok {
		d.h.ServeHTTP(w, r)
		return
	}

	h.ServeHTTP(w, r)
}

// Serve ...
func (d *DomMux) Serve(port string) error {
	if err := http.ListenAndServe(port, d); err != nil {
		return fmt.Errorf("dommux: failed while serving: %s", err)
	}

	return nil
}
