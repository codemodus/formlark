package main

import (
	"flag"
	"os"
	"time"

	"github.com/codemodus/formlark/internal/api"
	"github.com/codemodus/formlark/internal/dommux"
	"github.com/codemodus/formlark/internal/front"
	"github.com/codemodus/sigmon"
	"github.com/codemodus/vitals"
	"github.com/sirupsen/logrus"
)

type scopes struct {
	sm  string
	cnf string
	prf string
	srv string
}

func (s *scopes) String() string {
	return "----"
}

func main() {
	http := ":54541"
	profCPU := ""
	var statsCyc time.Duration
	profMem := ""

	log := logrus.New()
	scp := scopes{
		sm:  "sigmon",
		cnf: "conf",
		prf: "prof",
	}

	sm := sigmon.New(nil)
	log.Infof("%s: initialized", scp.sm)
	log.Infof("%s: set (ignoring)", scp.sm)

	sm.Run()
	log.Infof("%s: running", scp.sm)

	flag.StringVar(
		&http, "http", http,
		"port to listen on for http requests",
	)
	flag.StringVar(
		&profCPU, "prof-cpu", profCPU,
		"location to dump CPU profile",
	)
	flag.DurationVar(
		&statsCyc, "stats-cyc", statsCyc,
		"report memory stats with cycle length",
	)
	flag.StringVar(
		&profMem, "prof-mem", profMem,
		"location to dump memory profile",
	)

	flag.Parse()
	log.Infof("%s: flags processed", scp.cnf)

	stopCPUProf, err := vitals.StartCPUProfile(profCPU)
	if err != nil {
		log.Fatalf("%s: cpu failed: %s", scp.prf, err)
	}
	defer stopCPUProf()
	if profCPU != "" {
		log.Infof("%s: cpu > %q", scp.prf, profCPU)
	}

	if statsCyc > 0 {
		log.Infof("%s: mem stats > stdout", scp.prf)
	}
	c := vitals.MonitorMemoryStats(statsCyc)
	go func() {
		for statsCyc > 0 {
			log.Info(<-c)
		}
	}()

	h, err := dommux.New(
		dommux.WithDomainHandler("www.formlark.localhost", &front.Front{}),
		dommux.WithDomainHandler("api.formlark.localhost", &api.API{}),
	)
	if err != nil {
		log.Fatalf("%s: failed to initialize: %s", scp.srv, err)
	}

	sm.Set(func(ssm *sigmon.SignalMonitor) {
		os.Exit(0)
	})

	if err := h.Serve(http); err != nil {
		log.Errorf("%s: failed on serve: %s", scp.srv, err)
	}

	if profMem != "" {
		log.Infof("%s: heap > %q", scp.prf, profMem)
	}
	if err := vitals.WriteHeapProfile(profMem); err != nil {
		log.Fatalf("%s: heap failed: %s", scp.prf, err)
	}
}
