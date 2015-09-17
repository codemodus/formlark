package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

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

	d := n.newPageAuthed()
	d.URLLogin = "/" + n.su.conf.AdminPathPrefix + "/login"
	d.PageTitle = "Login"
	d.Footer.ColsDropdownFlag = false
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

func (n *node) adminOverviewHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		*PageAdmin
		User string
	}{
		n.newPageAdmin(),
		usr,
	}

	usrs := n.newUsers(5)
	if err := usrs.get(0); err != nil {
		fmt.Println(err)
	}
	fmt.Println(usrs)

	d.PageTitle = "Overview"
	n.ExecuteTemplate(w, "admin", d)
}

func (n *node) adminUsersHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		*PageAdmin
		Misc string
	}{
		n.newPageAdmin(),
		strconv.FormatInt(t, 10),
	}
	d.PageTitle = "Users"
	n.ExecuteTemplate(w, "admin/users", d)
}

func (n *node) adminSettingsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		*PageAdmin
		Misc string
	}{
		n.newPageAdmin(),
		strconv.FormatInt(t, 10),
	}
	d.PageTitle = "Settings"
	n.ExecuteTemplate(w, "admin/settings", d)
}
