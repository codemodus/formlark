package main

import (
	"bytes"
	"fmt"
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
	/*
		unreg -

		reg -
			store submitted form for later sending (encrypt)
			get email and send form
	*/

	si := 0
	if n.su.conf.FormPathPrefix != "" {
		si = 1
	}
	seg, err := parth.SegmentToString(r.URL.Path, si)
	if err != nil {
		http.Error(w, "cannot process path", 422)
		return
	}

	k, err := n.su.ds.dcbMrks.getBytes(seg)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if len(k) > 0 {
		u := n.newUser()
		u.ID = string(k)
		if err = u.get(); err != nil {
			http.Error(w, "user index found; cannot resolve user data - please contact the admin", 500)
			return
		}
		// DO STUFF
		rdr, err := n.su.ds.dcbAsts.get(u.PublicID)
		if err != nil {
			n.su.logs.Err.Println(err)
		}
		b := &bytes.Buffer{}
		b.ReadFrom(rdr)
		fmt.Fprintf(w, string(b.Bytes()))
		return

		if err := u.set(); err != nil {
			http.Error(w, "cannot persist to datastore", 500)
			return
		}
		// RESPOND
	}

	if strings.Contains(seg, "@") {
		a, err := mail.ParseAddress(seg)
		if err != nil {
			http.Error(w, "cannot process email", 422)
			return
		}
		u := n.newUser()
		u.Email = a.Address
		u.ID = n.getKey()
		u.PublicID = u.ID
		u.Confirm = &dtConfirm{}
		u.Confirm.Email = n.getConfirmHash()
		u.Confirm.Forms = make(map[string]string)
		u.Confirm.Forms[r.Referer()] = n.getConfirmHash()

		if err := u.set(); err != nil {
			http.Error(w, "cannot persist to datastore", 500)
			return
		}

		rdr, err := n.su.ds.dcbAsts.get(u.PublicID)
		if err != nil {
			n.su.logs.Err.Println(err)
		}
		b := &bytes.Buffer{}
		b.ReadFrom(rdr)
		fmt.Fprintf(w, string(b.Bytes()))
		return
	}
	/*
		rdr, err := n.su.ds.dcbAsts.get("test")
		if err != nil {
			n.su.logs.Err.Println(err)
		}
		b := &bytes.Buffer{}
		b.ReadFrom(rdr)
		fmt.Fprintf(w, string(b.Bytes()))
	*/
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
