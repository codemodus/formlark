package main

import (
	"html/template"
	"time"

	"github.com/codemodus/loggers"
	"golang.org/x/net/context"
)

const (
	EX__BASE       = 64 /* base value for error messages */
	EX_USAGE       = 64 /* command line usage error */
	EX_DATAERR     = 65 /* data format error */
	EX_NOINPUT     = 66 /* cannot open input */
	EX_NOUSER      = 67 /* addressee unknown */
	EX_NOHOST      = 68 /* host name unknown */
	EX_UNAVAILABLE = 69 /* service unavailable */
	EX_SOFTWARE    = 70 /* internal software error */
	EX_OSERR       = 71 /* system error (e.g., can't fork) */
	EX_OSFILE      = 72 /* critical OS file missing */
	EX_CANTCREAT   = 73 /* can't create (user) output file */
	EX_IOERR       = 74 /* input/output error */
	EX_TEMPFAIL    = 75 /* temp failure; user is invited to retry */
	EX_PROTOCOL    = 76 /* remote error in protocol */
	EX_NOPERM      = 77 /* permission denied */
	EX_CONFIG      = 78 /* configuration error */
)

type sysUtils struct {
	conf *conf
	ds   *dataStores
	logs *loggers.Loggers
	ts   *template.Template
}

type rCtxCmn int

const (
	keyPostHandlerFuncCtx rCtxCmn = iota
	keyReqStart
)

func SetReqStart(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, keyReqStart, t)
}

func GetReqStart(ctx context.Context) (time.Time, bool) {
	t, ok := ctx.Value(keyReqStart).(time.Time)
	return t, ok
}

// InitPHFC takes a context.Context and places a pointer to it within itself.
// This is useful for carrying data into the post ServeHTTPContext area of
// Handler wraps.  PHFC stands for Post HandlerFunc Context.
func InitPHFC(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyPostHandlerFuncCtx, &ctx)
}

// GetPHFC takes a context.Context and returns a pointer to the context.Context
// set in InitPHFC.
func GetPHFC(ctx context.Context) (*context.Context, bool) {
	cx, ok := ctx.Value(keyPostHandlerFuncCtx).(*context.Context)
	return cx, ok
}
