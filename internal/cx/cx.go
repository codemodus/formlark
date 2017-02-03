package cx

import (
	"context"
	"net/http"
)

type ctxKey int

// Key constants.
const (
	HTTPAuthKey ctxKey = iota
	HTTPTempAuthKey
)

// ReqSetHTTPAuth ...
func ReqSetHTTPAuth(r *http.Request, auth uint64) {
	nr := r.WithContext(SetHTTPAuth(r.Context(), auth))
	*r = *nr
}

// SetHTTPAuth ...
func SetHTTPAuth(ctx context.Context, auth uint64) context.Context {
	return context.WithValue(ctx, HTTPAuthKey, auth)
}

// HTTPAuth ...
func HTTPAuth(ctx context.Context) (uint64, bool) {
	auth, ok := ctx.Value(HTTPAuthKey).(uint64)
	return auth, ok
}

// ReqSetHTTPTempAuth ...
func ReqSetHTTPTempAuth(r *http.Request, auth uint64) {
	nr := r.WithContext(SetHTTPTempAuth(r.Context(), auth))
	*r = *nr
}

// SetHTTPTempAuth ...
func SetHTTPTempAuth(ctx context.Context, auth uint64) context.Context {
	return context.WithValue(ctx, HTTPTempAuthKey, auth)
}

// HTTPTempAuth ...
func HTTPTempAuth(ctx context.Context) (uint64, bool) {
	auth, ok := ctx.Value(HTTPTempAuthKey).(uint64)
	return auth, ok
}
