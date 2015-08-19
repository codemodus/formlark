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
		if err = u.validate(); err != nil {
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
