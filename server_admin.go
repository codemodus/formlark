package main

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

func (n *node) adminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s := n.sm.SessStart(w, r)
	usr, ok := s.Get("user")
	if !ok {
		http.Error(w, "bad session var", 500)
		return
	}
	u := usr.(string)
	s.Set("test", time.Now().Unix())

	d := struct {
		*Page
		User string
	}{
		n.newPage(),
		u,
	}
	err := n.su.ts.ExecuteTemplate(w, "admin/index.html", d)
	if err != nil {
		http.Error(w, "template failed - please contact the site admin", 500)
		return
	}
	return
}

func (n *node) adminLoginGetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s := n.sm.SessStart(w, r)
	usr, ok := s.Get("user")
	if ok && usr.(string) != "" {
		http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix, 302)
		return
	}

	p := n.newPage()
	err := n.su.ts.ExecuteTemplate(w, "admin/login.html", p)
	if err != nil {
		http.Error(w, "template failed - please contact the site admin", 500)
		return
	}
	return
}

func (n *node) adminLoginPostHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "cannot parse form", 422)
		return
	}
	usr := r.Form.Get("user")
	pass := r.Form.Get("pass")
	if usr == n.su.conf.AdminUser && pass == n.su.conf.AdminPass {
		s := n.sm.SessStart(w, r)
		s.Set("user", usr)
		s.Set("test", time.Now().Unix())

		http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix, 303)
		return
	}

	http.Error(w, "unauthorized", 401)
	return
}

func (n *node) adminTestHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s := n.sm.SessStart(w, r)
	t, ok := s.Get("test")
	if !ok {
		http.Error(w, "bad session var", 500)
		return
	}

	d := struct {
		*Page
		Misc string
	}{
		n.newPage(),
		strconv.FormatInt(t.(int64), 10),
	}
	err := n.su.ts.ExecuteTemplate(w, "admin/test.html", d)
	if err != nil {
		http.Error(w, "template failed - please contact the site admin", 500)
		return
	}
	return
}
