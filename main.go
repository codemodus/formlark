package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codemodus/config"
	"github.com/codemodus/loggers"
	"github.com/codemodus/sigmon"
)

type sysUtils struct {
	conf *conf
	ds   *dataStores
	logs *loggers.Loggers
	ts   *Templates
}

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
		os.Exit(1)
	}

	lOpts := loggers.NewBypassedOptions()
	lOpts.BypassSys = false
	lOpts.BypassErr = false
	lOpts.SysToStdout = true
	lOpts.ErrToStderr = true
	if su.logs, err = loggers.New(lOpts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if su.ds.dcbsRsrcs, err = getDataCacheLocal("test.db"); err != nil {
		su.logs.Err.Println(err)
		os.Exit(1)
	}
	if su.ds.dcbUsers, err = su.ds.dcbsRsrcs.getBucket("users"); err != nil {
		su.logs.Err.Println(err)
		os.Exit(1)
	}
	if su.ds.dcbIndUsers, err = su.ds.dcbsRsrcs.getBucket("index-users"); err != nil {
		su.logs.Err.Println(err)
		os.Exit(1)
	}
	if su.ds.dcbIndCnfrm, err = su.ds.dcbsRsrcs.getBucket("index-confirmation"); err != nil {
		su.logs.Err.Println(err)
		os.Exit(1)
	}
	if su.ds.dcbPosts, err = su.ds.dcbsRsrcs.getBucket("posts"); err != nil {
		su.logs.Err.Println(err)
		os.Exit(1)
	}

	su.ts = NewTemplates("", "")
	su.ts.ParseDir("front/templates")

	cl := newCluster(su)
	cl.Configure(false)
	cl.Run()
	sigMon.Set(cl.signal)
	cl.Wait()
	sigMon.Stop()
}
