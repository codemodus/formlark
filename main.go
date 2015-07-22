package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codemodus/config"
	"github.com/codemodus/loggers"
	"github.com/codemodus/sigmon"
)

func main() {
	sigMon := sigmon.New(nil)
	sigMon.Run()

	su := &sysUtils{
		conf: &conf{},
		ds:   &dataStores{},
	}
	err := config.Init(su.conf, path.Join(config.DefaultDir, config.DefaultFilename))
	if err != nil {
		fmt.Println(err)
		os.Exit(EX_CONFIG)
	}

	lOpts := loggers.NewBypassedOptions()
	lOpts.BypassSys = false
	lOpts.BypassErr = false
	lOpts.SysToStdout = true
	lOpts.ErrToStderr = true
	if su.logs, err = loggers.New(lOpts); err != nil {
		fmt.Println(err)
		os.Exit(EX_CANTCREAT)
	}

	if su.ds.dcbsRsrcs, err = getDataCacheLocal("test.db"); err != nil {
		su.logs.Err.Println(err)
		os.Exit(EX_CANTCREAT)
	}
	if su.ds.dcbAsts, err = su.ds.dcbsRsrcs.getBucket("assets"); err != nil {
		su.logs.Err.Println(err)
		os.Exit(EX_CANTCREAT)
	}

	cl := newCluster(su)
	cl.Configure(false)
	cl.Run()
	sigMon.Set(cl.signal)
	cl.Wait()
	sigMon.Stop()
}
