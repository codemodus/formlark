package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/codemodus/formlark/internal/datatypes"
)

type user struct {
	*boltItem
	*datatypes.User
}

func (n *node) newUser() *user {
	u := &user{
		User: datatypes.NewUser(),
		boltItem: &boltItem{
			BB:  n.u.ds.dcbUsers,
			BBI: n.u.ds.dcbIndUsers,
		},
	}
	return u
}

func (u *user) affixID() error {
	if u.ID != "" {
		if err := u.BB.find(u.ID); err == nil {
			return nil
		}
	}
	if u.Email != "" {
		k, err := u.BBI.getBytes(u.Email)
		if err != nil {
			return err
		}
		if len(k) > 0 {
			u.ID = string(k)
			return nil
		}
	}
	return errors.New("not found")
}

func (u *user) set() error {
	err := u.boltItem.set()
	if err != nil {
		return err
	}
	if err = u.BBI.setBytes(u.Email, []byte(u.ID)); err != nil {
		return err
	}
	return nil
}

type users struct {
	BI *boltItem
	s  []*user
}

func (n *node) newUsers(count int) *users {
	u := &users{
		s: make([]*user, count),
		BI: &boltItem{
			BB:  n.u.ds.dcbUsers,
			BBI: n.u.ds.dcbIndUsers,
		},
	}
	return u
}

func (us *users) get(start int) error {
	m, err := us.BI.BB.getManyBytes(start, cap(us.s))
	if err != nil {
		return err
	}

	ct := -1
	for k, v := range m {
		ct++
		br := bytes.NewReader(v)
		dec := gob.NewDecoder(br)
		tmp := &user{
			User: datatypes.NewUser(),
			boltItem: &boltItem{
				BB:  us.BI.BB,
				BBI: us.BI.BB,
			},
		}
		if err := dec.Decode(tmp); err != nil {
			return err
		}
		us.s[ct] = tmp
		us.s[ct].ID = k
	}
	return nil
}

type post struct {
	*datatypes.Post
}

func (n *node) newPost() *post {
	return &post{
		Post: datatypes.NewPost(),
	}
}

func (p *post) processForm(r *http.Request) error {
	p.Date = time.Now()
	p.Subject = "Message submitted through " + r.Referer()
	if r.Form.Get("email") != "" {
		p.Replyto = r.Form.Get("email")
	}
	for k, v := range r.Form {
		if k == "_replyto" {
			p.Replyto = v[0]
			continue
		}
		if k == "_next" {
			p.Next = v[0]
			continue
		}
		if k == "_subject" {
			p.Subject = v[0]
			continue
		}
		if k == "_cc" {
			p.CC = v
			continue
		}
		p.Content[k] = v
	}
	return nil
}

type posts struct {
	*boltItem
	S []*post
}

func (n *node) newPosts(i string) *posts {
	return &posts{
		S: make([]*post, 0),
		boltItem: &boltItem{
			ID: i, BB: n.u.ds.dcbPosts,
		},
	}
}
