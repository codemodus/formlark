package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codemodus/chain"
	"github.com/codemodus/loggers"
	"golang.org/x/net/context"
)

func (n *node) reco(next chain.Handler) chain.Handler {
	return chain.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTPContext(ctx, w, r)
	})
}

func (n *node) initReq(next chain.Handler) chain.Handler {
	return chain.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = n.SetReqStart(ctx, time.Now())
		ctx = n.InitPHFC(ctx)
		next.ServeHTTPContext(ctx, w, r)
	})
}

func (n *node) log(next chain.Handler) chain.Handler {
	return chain.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		next.ServeHTTPContext(ctx, w, r)

		t2 := time.Now()
		t1, _ := n.GetReqStart(ctx)
		dur := t2.Sub(t1)
		str := fmt.Sprintf(loggers.CLF, r.RemoteAddr, r.Host,
			r.Method, r.URL.String(), dur)

		go func(s *node, toLog string) {
			//s.su.logs.Dbg.Will(s.su.logs).Print(toLog)
		}(n, str)

		/*pc, _ := chain.GetPHFC(ctx)
		tx1, _ := startTimeFromCtx(*pc)*/
	})
}

func (n *node) sess(next chain.Handler) chain.Handler {
	return chain.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		s, err := n.sm.SessStart(w, r)
		if err != nil {
			fmt.Println("ouch")
			// TODO
		}

		if r.URL.Path[len(r.URL.Path)-7:] == "/logout" {
			n.sm.SessStop(w, r)
			if r.URL.Path[1:len(n.u.conf.AdminPathPrefix)+1] == n.u.conf.AdminPathPrefix {
				http.Redirect(w, r, "/"+n.u.conf.AdminPathPrefix+"/login", 302)
				return
			}
			http.Redirect(w, r, "/login", 302)
			return
		}

		usr, ok := s.Get("user").(string)
		if !ok || usr == "" {
			s.Set("prevReq", r.URL.Path)
			if r.URL.Path[1:len(n.u.conf.AdminPathPrefix)+1] == n.u.conf.AdminPathPrefix {
				http.Redirect(w, r, "/"+n.u.conf.AdminPathPrefix+"/login", 302)
				return
			}
			http.Redirect(w, r, "/login", 302)
			return
		}

		ctx = n.SetSess(ctx, s)
		next.ServeHTTPContext(ctx, w, r)
	})
}
