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
	pvr := sessmgr.NewProviderRegistry()
	pvd := sessmgr.NewVolatileProvider()
	pvr.Register("test", pvd)
	sm, err := sessmgr.New(pvr, "test", "cook-e", 45)
	if err != nil {
		panic(err.Error())
	}

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

	m.Post(path.Join("/"+n.su.conf.FormPathPrefix+"/*x"), c.EndFn(n.postHandler))

	mAdm := m.Group("/" + n.su.conf.AdminPathPrefix)
	mAdm.Get("/", sc.EndFn(n.adminGetHandler))
	mAdm.Get("/login", c.EndFn(n.adminLoginGetHandler))
	mAdm.Post("/login", c.EndFn(n.adminLoginPostHandler))
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
		ctx = SetReqStart(ctx, time.Now())
		ctx = InitPHFC(ctx)
		next.ServeHTTPContext(ctx, w, r)
	})
}

func (n *node) log(next chain.Handler) chain.Handler {
	return chain.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		next.ServeHTTPContext(ctx, w, r)

		t2 := time.Now()
		t1, _ := GetReqStart(ctx)
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
		s := n.sm.SessStart(w, r)
		logged := s.Get("logged")
		if logged == nil || logged.(bool) == false {
			http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix+"/login", 302)
			return
		}
		/*lo := r.URL.Query().Get("logout")
		if lo != "" {
			//http.Error(w, "Logged out", 401)
			//return
			r.SetBasicAuth("", "")
		}*/
		next.ServeHTTPContext(ctx, w, r)
	})
}

func (n *node) NotFound(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	return
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
