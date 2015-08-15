package main

import (
	"net/http"

	"golang.org/x/net/context"
)

func (n *node) adminGetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s := n.sm.SessionStart(w, r)
	usr := s.Get("user")
	if usr == nil || usr.(string) == "" {
		http.Redirect(w, r, "/", 302)
		return
	}

	d := struct {
		*Page
		User string
	}{
		n.newPage(),
		usr.(string),
	}
	err := n.su.ts.ExecuteTemplate(w, "admin/index.html", d)
	if err != nil {
		http.Error(w, "template failed - please contact the site admin", 500)
		return
	}
	return
}

func (n *node) adminLoginGetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	p := n.newPage()
	err := n.su.ts.ExecuteTemplate(w, "admin/login.html", p)
	if err != nil {
		http.Error(w, "template failed - please contact the site admin", 500)
		return
	}
}

func (n *node) adminLoginPostHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "cannot parse form", 422)
		return
	}
	usr := r.Form.Get("user")
	pass := r.Form.Get("pass")
	if usr == n.su.conf.AdminUser && pass == n.su.conf.AdminPass {
		s := n.sm.SessionStart(w, r)
		s.Set("logged", true)
		s.Set("user", usr)

		http.Redirect(w, r, "/"+n.su.conf.AdminPathPrefix, 303)
		return
	}

	http.Error(w, "unauthorized", 401)
	return
}
