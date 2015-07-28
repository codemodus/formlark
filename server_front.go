package main

import (
	"hash/fnv"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/codemodus/parth"
	"golang.org/x/net/context"
)

func (n *node) postHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	k, err := n.su.ds.dcbMrks.getBytes(seg)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// If key is found...
	if len(k) > 0 {
		// Get user by key.
		u := n.newUser()
		u.ID = string(k)
		if err = u.get(); err != nil {
			http.Error(w, "user index found; cannot resolve user data - please contact the admin", 500)
			return
		}

		// Prep request form.
		if err = r.ParseForm(); err != nil {
			http.Error(w, "cannot process form", 422)
			return
		}

		// Get posts if exists or create new.
		ps := n.newPosts()
		ps.ID = u.ID
		if err = n.su.ds.dcbPosts.find(u.ID); err != nil {
			if err = ps.get(); err != nil {
				http.Error(w, "posts index found; cannot resolve posts data - please contact the admin", 500)
				return
			}
		}
		// Process form, validate, add to posts, and persist.
		p := n.newPost()
		if err = p.processForm(r); err != nil {
			http.Error(w, "cannot process form", 422)
			return
		}
		// TODO: Add mail validation to posts (replyto, cc).
		// TODO: Add url validation to posts (next).
		// TODO: Validate.
		ps.S = append(ps.S, p)
		if err = ps.set(); err != nil {
			http.Error(w, "cannot persist post to datastore", 500)
			return
		}

		http.Redirect(w, r, p.Next, 303)
		return
	}

	// If key is not found and segment might be an email...
	if strings.Contains(seg, "@") {
		// TODO: Move/Add mail validation to user.
		a, err := mail.ParseAddress(seg)
		if err != nil {
			http.Error(w, "cannot process email", 422)
			return
		}
		// Create new user.
		u := n.newUser()
		u.Email = a.Address
		u.ID = n.getKey()
		u.Confirm = &dtConfirm{}
		u.Confirm.Email = n.getConfirmHash()
		u.Confirm.Forms = make(map[string]string)
		// Check for empty referer.
		u.Confirm.Forms[r.Referer()] = n.getConfirmHash()

		// TODO: Validate.

		if err := u.set(); err != nil {
			http.Error(w, "cannot persist user to datastore", 500)
			return
		}

		// TODO: Add confirmation email to "spool".

		// TODO: Place form processing into function and call here.
		return
	}
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

func (n *node) processConfirm(r *http.Request, u *user) {

}
