package main

import (
	"flag"
	"time"

	"github.com/codemodus/sigmon"
	"github.com/codemodus/vitals"
	"github.com/sirupsen/logrus"
)

type scopes struct {
	sm  string
	cnf string
	prf string
}

func (s *scopes) String() string {
	return "----"
}

func main() {
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

	// TODO: init

	sm.Set(func(ssm *sigmon.SignalMonitor) {
		// TODO: sigs
	})

	// TODO: run
	time.Sleep(time.Second * 18)

	if profMem != "" {
		log.Infof("%s: heap > %q", scp.prf, profMem)
	}
	if err := vitals.WriteHeapProfile(profMem); err != nil {
		log.Fatalf("%s: heap failed: %s", scp.prf, err)
	}
}
