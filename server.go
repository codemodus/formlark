package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/codemodus/chain"
	"github.com/codemodus/formlark/internal/sessmgr"
	"github.com/codemodus/httpcluster"
	"github.com/codemodus/loggers"
	"github.com/codemodus/mixmux"
	"github.com/codemodus/sigmon"
	"golang.org/x/net/context"
)

type node struct {
	*httpcluster.Node
	su *sysUtils
	sm *sessmgr.Manager
}

type cluster struct {
	*httpcluster.Cluster
	su *sysUtils
}

func newCluster(su *sysUtils) *cluster {
	return &cluster{
		Cluster: &httpcluster.Cluster{}, su: su,
	}
}

func (cl *cluster) Configure(linkage bool) {
	p := sessmgr.NewVolatileProvider()
	sm := sessmgr.New("cook-e", 45, p)

	n := &node{
		su: cl.su, sm: sm,
		Node: &httpcluster.Node{
			Timeout: time.Second * 5, Addr: cl.su.conf.ServerPort,
		},
	}
	n.ErrorLog = n.su.logs.Err.Logger
	n.Handler = n.setupMux()

	cl.AddNode(n.Node)
}

func (n *node) setupMux() *mixmux.TreeMux {
	c := chain.New(n.reco, n.initReq, n.log, chain.Convert(n.Node.Wedge))
	sc := c.Append(n.auth)
	m := mixmux.NewTreeMux()

	m.Get("/assets/public/*x", c.EndFn(n.assetsHandler))
	m.Get("/assets/protected/*x", sc.EndFn(n.assetsHandler))
	m.Post(path.Join("/"+n.su.conf.FormPathPrefix+"/*x"), c.EndFn(n.postHandler))

	mAdm := m.Group("/" + n.su.conf.AdminPathPrefix)
	mAdm.Get("/", sc.EndFn(n.adminHandler))
	mAdm.Get("/login", c.EndFn(n.adminLoginGetHandler))
	mAdm.Post("/login", c.EndFn(n.adminLoginPostHandler))
	mAdm.Get("/test", sc.EndFn(n.adminTestHandler))
	mAdm.Get("/*x", c.EndFn(n.NotFound))
	return m
}

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

func (n *node) auth(next chain.Handler) chain.Handler {
	return chain.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		s, err := n.sm.SessStart(w, r)
		if err != nil {
			fmt.Println("ouch")
			// TODO
		}
		usr, ok := s.Get("user").(string)
		if !ok || usr == "" {
			http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix+"/login", 302)
			return
		}

		out := r.URL.Query().Get("logout")
		if out != "" {
			n.sm.SessStop(w, r)
			http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix+"/login", 302)
			return
		}

		ctx = n.SetSess(ctx, s)
		next.ServeHTTPContext(ctx, w, r)
	})
}

func (n *node) NotFound(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (cl *cluster) signal(sm *sigmon.SignalMonitor) {
	switch sm.Sig() {
	case sigmon.SIGHUP:
		cl.Stop(false)
		cl.Run()
	case sigmon.SIGINT, sigmon.SIGTERM:
		cl.Stop(true)
	case sigmon.SIGUSR1, sigmon.SIGUSR2:
		//
	}
}

func (n *node) assetsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p[0] == '/' {
		p = p[1:]
	}
	http.ServeFile(w, r, p)
}
