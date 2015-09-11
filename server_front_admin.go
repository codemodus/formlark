package main

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

func (n *node) adminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s, ok := n.GetSess(ctx)
	if !ok {
		http.Error(w, "session not found in context", 500)
		return
	}
	usr, ok := s.Get("user").(string)
	if !ok {
		http.Error(w, "bad session var", 500)
		return
	}
	s.Set("test", time.Now().Unix())

	d := struct {
		*Page
		User string
	}{
		n.newPage(),
		usr,
	}
	n.ExecuteTemplate(w, "admin", d)
}

func (n *node) adminLoginGetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s, err := n.sm.SessStart(w, r)
	if err != nil {
		//
	}
	usr, ok := s.Get("user").(string)
	if ok && usr != "" {
		http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix, 302)
		return
	}

	d := n.newPage()
	d.NavDrawer.NavCommonItems[0].Name = "TATA!!!"
	n.ExecuteTemplate(w, "admin/login", d)
}

func (n *node) adminLoginPostHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "cannot parse form", 422)
		return
	}
	usr := r.Form.Get("user")
	pass := r.Form.Get("pass")
	if usr == n.su.conf.AdminUser && pass == n.su.conf.AdminPass {
		s, err := n.sm.SessStart(w, r)
		if err != nil {
			// TODO
		}
		s.Set("user", usr)
		s.Set("test", time.Now().Unix())

		p, ok := s.Get("prevReq").(string)
		if ok && p != "" {
			http.Redirect(w, r, p, 302)
			return
		}

		http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix, 303)
		return
	}

	http.Error(w, "unauthorized", 401)
}

func (n *node) adminTestHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s, ok := n.GetSess(ctx)
	if !ok {
		http.Error(w, "session not found in context", 500)
		return
	}
	t, ok := s.Get("test").(int64)
	if !ok {
		http.Error(w, "bad session var", 500)
		return
	}

	d := struct {
		*Page
		Misc string
	}{
		n.newPage(),
		strconv.FormatInt(t, 10),
	}
	n.ExecuteTemplate(w, "admin/test", d)
}
