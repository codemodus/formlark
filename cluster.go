package main

import (
	"time"

	"github.com/codemodus/httpcluster"
	"github.com/codemodus/sessctrl"
	"github.com/codemodus/sigmon"
)

type cluster struct {
	*httpcluster.Cluster
	u *utils
}

func newCluster(u *utils) *cluster {
	return &cluster{
		Cluster: &httpcluster.Cluster{}, u: u,
	}
}

func (cl *cluster) Init() {
	p := sessctrl.NewVolatileProvider()
	sc := sessctrl.NewCookieController("cook-e", 90, p)

	n := &node{
		u: cl.u, sc: sc,
		Node: &httpcluster.Node{
			Timeout: time.Second * 5, Addr: cl.u.conf.ServerPort,
		},
	}
	n.ErrorLog = n.u.logs.Err.Logger
	n.Handler = n.setupMux()

	cl.AddNode(n.Node)
}

func (cl *cluster) signal(sm *sigmon.SignalMonitor) {
	switch sm.Sig() {
	case sigmon.SIGHUP:
		cl.Restart(nil)
	case sigmon.SIGINT, sigmon.SIGTERM:
		cl.Stop()
	case sigmon.SIGUSR1, sigmon.SIGUSR2:
		//
	}
}
