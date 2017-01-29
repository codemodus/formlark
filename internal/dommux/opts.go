package dommux

import (
	"fmt"
	"net/http"
)

type options struct {
	m map[string]http.Handler
	h http.Handler
}

// Option ...
type Option func(*options) error

// WithDefaultHandler ...
func WithDefaultHandler(handler http.Handler) Option {
	return func(opts *options) error {
		if handler == nil {
			return fmt.Errorf("default handler opt: must provide valid handler")
		}

		opts.h = handler

		return nil
	}
}

// WithDomainHandler ...
func WithDomainHandler(domain string, handler http.Handler) Option {
	return func(opts *options) error {
		if domain == "" {
			return fmt.Errorf("domain handler opt: must provide valid domain")
		}
		if handler == nil {
			return fmt.Errorf("domain handler opt: must provide valid handler")
		}

		opts.m[domain] = handler

		return nil
	}
}
