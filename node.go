package main

import (
	"errors"
	"hash/fnv"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/codemodus/chain"
	"github.com/codemodus/formlark/internal/sessmgr"
	"github.com/codemodus/httpcluster"
	"github.com/codemodus/mixmux"
	"github.com/codemodus/parth"
)

type node struct {
	*httpcluster.Node
	u  *utils
	sm *sessmgr.Manager
}

func (n *node) setupMux() *mixmux.TreeMux {
	c := chain.New(n.reco, n.initReq, n.log)
	s := c.Append(n.sess)
	m := mixmux.NewTreeMux()

	m.Get("/favicon.ico", c.EndFn(n.iconHandler))

	m.Get("/", c.EndFn(n.anonIndexHandler))
	m.Get("/login", c.EndFn(n.authedLoginGetHandler))
	m.Post("/login", c.EndFn(n.authedLoginPostHandler))
	m.Get("/logout", s.EndFn(n.NotFound))

	m.Get("/overview", s.EndFn(n.authedOverviewHandler))
	m.Get("/settings", s.EndFn(n.authedSettingsHandler))

	m.Get("/assets/*x", c.EndFn(n.assetsHandler))
	m.Get("/jspm_packages/*x", c.EndFn(n.assetsFlexHandler))
	m.Get("/app/*x", s.EndFn(n.assetsFlexHandler))

	m.Post(path.Join("/"+n.u.conf.FormPathPrefix+"/*x"), c.EndFn(n.anonPostHandler))

	mA := m.Group("/" + n.u.conf.AdminPathPrefix)
	mA.Get("/", s.EndFn(n.adminOverviewHandler))
	mA.Get("/login", c.EndFn(n.adminLoginGetHandler))
	mA.Post("/login", c.EndFn(n.adminLoginPostHandler))
	mA.Get("/logout", s.EndFn(n.NotFound))

	mA.Get("/overview", s.EndFn(n.adminOverviewHandler))
	mA.Get("/users", s.EndFn(n.adminUsersHandler))
	mA.Get("/settings", s.EndFn(n.adminSettingsHandler))
	mA.Get("/backup", s.EndFn(n.backupHandleFunc))
	mA.Get("/app/*x", s.EndFn(n.assetsFlexHandler))

	mA.Get("/assets/*x", s.EndFn(n.assetsHandler))
	mA.Get("/jspm_packages/*x", s.EndFn(n.assetsFlexHandler))

	mA.Get("/*x", c.EndFn(n.NotFound))
	return m
}

func (n *node) getReferer(r string) (string, error) {
	ref, err := url.Parse(r)
	if err != nil || ref == nil {
		return "", errors.New("error parsing referer: " + err.Error())
	}
	return ref.String(), nil
}

func (n *node) getIndexSegment(s string) (string, error) {
	si := 0
	if n.u.conf.FormPathPrefix != "" {
		si = 1
	}
	seg, err := parth.SegmentToString(s, si)
	if err != nil {
		return "", err
	}
	return seg, nil
}

func (n *node) getKey() string {
	t := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	h := fnv.New64a()
	h.Write(t)
	s := strconv.FormatUint(h.Sum64(), 10)
	return s
}

func (n *node) getConfirmHash() string {
	t := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	h := fnv.New64a()
	h.Write(t)
	s := strconv.FormatUint(h.Sum64(), 10) +
		"_" + strconv.FormatInt(time.Now().Unix(), 10)
	return s
}
