package main

import (
	"errors"
	"hash/fnv"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/codemodus/parth"
	"golang.org/x/net/context"
)

func (n *node) loginGetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s, err := n.sm.SessStart(w, r)
	if err != nil {
		//
	}
	usr, ok := s.Get("user").(string)
	if ok && usr != "" {
		http.Redirect(w, r, "/overview", 302)
		return
	}

	d := n.newPagePublic()
	d.PageTitle = "Login"
	n.ExecuteTemplate(w, "login", d)
}

func (n *node) loginPostHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func (n *node) overviewHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		*PagePublic
		User string
	}{
		n.newPagePublic(),
		usr,
	}
	d.PageTitle = "Overview"
	n.ExecuteTemplate(w, "index", d)
}

func (n *node) settingsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		*PagePublic
		User string
	}{
		n.newPagePublic(),
		usr,
	}
	d.PageTitle = "Settings"
	n.ExecuteTemplate(w, "index", d)
}

func (n *node) postHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	rf, err := n.getReferer(r.Referer())
	if err != nil {
		http.Error(w, "referer must be parsable", 400)
		return
	}
	seg, err := n.getIndexSegment(r.URL.Path)
	if err != nil {
		http.Error(w, "cannot process path", 422)
		return
	}

	u := n.newUser(seg)
	ok, err := u.get()
	if err != nil {
		http.Error(w, "cannot access user data in datastore", 500)
		return
	}
	if !ok {
		if u.Email == "" {
			http.NotFound(w, r)
			return
		}
		u.ID = n.getKey()
		if err = u.Validate(); err != nil {
			http.Error(w, "user data invalid: "+err.Error(), 422)
			return
		}

		if err = u.set(); err != nil {
			http.Error(w, "cannot persist user to datastore", 500)
			return
		}
	}

	// Check and handle form confirmation status.
	fConfirm, ok := u.Confirm.Forms[rf]
	if !ok || fConfirm != "" {
		if err = n.su.ds.dcbIndCnfrm.find(u.ID); err != nil {
			if err = n.su.ds.dcbIndCnfrm.setBytes(u.ID, []byte("")); err != nil {
				http.Error(w, "cannot persist confirmation index to datastore", 500)
				return
			}
		}
		if !ok {
			u.Confirm.Forms[rf] = n.getConfirmHash()
			// TODO: send message with all needed confirmations
			// TODO: persist user
		}
		http.Redirect(w, r, n.su.conf.ServerProtocol+n.su.conf.ServerDomain+"/unconfirmed", 303)
		return
	}

	ps := n.newPosts(u.ID)
	ok, err = ps.get()
	if err != nil {
		http.Error(w, "cannot access posts data in datastore", 500)
		return
	}

	// Prep request form.
	if err = r.ParseForm(); err != nil {
		http.Error(w, "cannot parse form", 422)
		return
	}

	// Process form, validate, add to posts, and persist.
	p := n.newPost()
	if err = p.processForm(r); err != nil {
		http.Error(w, "cannot process form", 422)
		return
	}
	if err = p.Validate(); err != nil {
		http.Error(w, "post data invalid: "+err.Error(), 422)
		return
	}
	ps.S = append(ps.S, p)

	if err = ps.set(); err != nil {
		http.Error(w, "cannot persist post to datastore", 500)
		return
	}

	http.Redirect(w, r, p.Next, 303)
	return
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
	if n.su.conf.FormPathPrefix != "" {
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
