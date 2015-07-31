package main

import (
	"hash/fnv"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/codemodus/parth"
	"golang.org/x/net/context"
)

func (n *node) postHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ref, err := url.Parse(r.Referer())
	if err != nil || ref == nil {
		http.Error(w, "referer must be parsable", 400)
		return
	}
	rf := ref.String()

	// Set segment index based on existence of form path prefix and get segment.
	si := 0
	if n.su.conf.FormPathPrefix != "" {
		si = 1
	}
	seg, err := parth.SegmentToString(r.URL.Path, si)
	if err != nil {
		http.Error(w, "cannot process path", 422)
		return
	}

	// Search for user key by email/id.
	k, err := n.su.ds.dcbIndUsers.getBytes(seg)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	u := n.newUser()

	// If key is found...
	if len(k) > 0 {
		u.ID = string(k)
		if err = u.get(); err != nil {
			http.Error(w, "cannot resolve user data - please contact the admin", 500)
			return
		}
	}

	// If key is not found and segment might be an email...
	if strings.Contains(seg, "@") {
		u.Email = seg
		u.ID = n.getKey()

		if err = u.validate(); err != nil {
			http.Error(w, "user data invalid: "+err.Error(), 422)
			return
		}

		if err := u.set(); err != nil {
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
		}
		http.Redirect(w, r, n.su.conf.ServerProtocol+n.su.conf.ServerDomain+"/unconfirmed", 303)
		return
	}

	ps := n.newPosts()
	ps.ID = u.ID
	if err = n.su.ds.dcbPosts.find(u.ID); err == nil {
		if err = ps.get(); err != nil {
			http.Error(w, "cannot resolve posts data - please contact the admin", 500)
			return
		}
	}

	// Prep request form.
	if err = r.ParseForm(); err != nil {
		http.Error(w, "cannot process form", 422)
		return
	}

	// Process form, validate, add to posts, and persist.
	p := n.newPost()
	if err = p.processForm(r); err != nil {
		http.Error(w, "cannot process form", 422)
		return
	}
	if err = p.validate(); err != nil {
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

func (n *node) nilHandler(w http.ResponseWriter, r *http.Request) {}

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
