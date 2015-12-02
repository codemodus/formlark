package main

import (
	"time"

	"github.com/codemodus/sessctrl"
	"golang.org/x/net/context"
)

type nodeCtxCommon int

const (
	nodeCtxKeyPHFC nodeCtxCommon = iota
	nodeCtxKeyReqStart
	nodeCtxKeySess
)

func (n *node) InitPHFC(ctx context.Context) context.Context {
	return context.WithValue(ctx, nodeCtxKeyPHFC, &ctx)
}

func (n *node) GetPHFC(ctx context.Context) (*context.Context, bool) {
	cx, ok := ctx.Value(nodeCtxKeyPHFC).(*context.Context)
	return cx, ok
}

func (n *node) SetReqStart(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, nodeCtxKeyReqStart, t)
}

func (n *node) GetReqStart(ctx context.Context) (time.Time, bool) {
	t, ok := ctx.Value(nodeCtxKeyReqStart).(time.Time)
	return t, ok
}

func (n *node) SetSess(ctx context.Context, s *sessctrl.Session) context.Context {
	return context.WithValue(ctx, nodeCtxKeySess, s)
}

func (n *node) GetSess(ctx context.Context) (*sessctrl.Session, bool) {
	s, ok := ctx.Value(nodeCtxKeySess).(*sessctrl.Session)
	return s, ok
}
