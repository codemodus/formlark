package main

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func (n *node) authedLoginGetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s, err := n.sm.SessStart(w, r)
	if err != nil {
		//
	}
	usr, ok := s.Get("user").(string)
	if ok && usr != "" {
		http.Redirect(w, r, "/overview", 302)
		return
	}

	d := n.newPageAnon()
	d.PageTitle = "Login"
	n.ExecuteTemplate(w, "login", d)
}

func (n *node) authedLoginPostHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

		http.Redirect(w, r, "/overview", 303)
		return
	}

	http.Error(w, "unauthorized", 401)
}

func (n *node) authedOverviewHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		*PageAuthed
		User string
	}{
		n.newPageAuthed(),
		usr,
	}
	d.PageTitle = "Overview"
	n.ExecuteTemplate(w, "index", d)
}

func (n *node) authedSettingsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		*PageAuthed
		User string
	}{
		n.newPageAuthed(),
		usr,
	}
	d.PageTitle = "Settings"
	n.ExecuteTemplate(w, "index", d)
}
