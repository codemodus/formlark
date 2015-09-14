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
	s := c.Append(n.sess)
	m := mixmux.NewTreeMux()

	m.Get("/", c.EndFn(n.anonIndexHandler))
	m.Get("/login", c.EndFn(n.authedLoginGetHandler))
	m.Post("/login", c.EndFn(n.authedLoginPostHandler))
	m.Get("/logout", s.EndFn(n.NotFound))

	m.Get("/overview", s.EndFn(n.authedOverviewHandler))
	m.Get("/settings", s.EndFn(n.authedSettingsHandler))

	m.Get("/assets/public/*x", c.EndFn(n.assetsHandler))
	m.Get("/assets/protected/*x", s.EndFn(n.assetsHandler))

	m.Post(path.Join("/"+n.su.conf.FormPathPrefix+"/*x"), c.EndFn(n.anonPostHandler))

	mA := m.Group("/" + n.su.conf.AdminPathPrefix)
	mA.Get("/", s.EndFn(n.adminOverviewHandler))
	mA.Get("/login", c.EndFn(n.adminLoginGetHandler))
	mA.Post("/login", c.EndFn(n.adminLoginPostHandler))
	mA.Get("/logout", s.EndFn(n.NotFound))

	mA.Get("/overview", s.EndFn(n.adminOverviewHandler))
	mA.Get("/users", s.EndFn(n.adminUsersHandler))
	mA.Get("/settings", s.EndFn(n.adminSettingsHandler))

	mA.Get("/*x", c.EndFn(n.NotFound))
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

func (n *node) sess(next chain.Handler) chain.Handler {
	return chain.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		s, err := n.sm.SessStart(w, r)
		if err != nil {
			fmt.Println("ouch")
			// TODO
		}

		if r.URL.Path[len(r.URL.Path)-7:] == "/logout" {
			n.sm.SessStop(w, r)
			if r.URL.Path[1:len(n.su.conf.AdminPathPrefix)+1] == n.su.conf.AdminPathPrefix {
				http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix+"/login", 302)
				return
			}
			http.Redirect(w, r, "/login", 302)
			return
		}

		usr, ok := s.Get("user").(string)
		if !ok || usr == "" {
			s.Set("prevReq", r.URL.Path)
			if r.URL.Path[1:len(n.su.conf.AdminPathPrefix)+1] == n.su.conf.AdminPathPrefix {
				http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix+"/login", 302)
				return
			}
			http.Redirect(w, r, "/login", 302)
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
