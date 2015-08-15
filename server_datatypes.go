package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

func (n *node) newPage() *Page {
	p := NewPage()
	p.URLLogin = "/" + n.su.conf.AdminPathPrefix + "/login"
	return p
}

type user struct {
	*boltItem
	*dtUser
}

func (n *node) newUser(i string) *user {
	u := &user{
		dtUser: newUser(),
		boltItem: &boltItem{
			DS: n.su.ds,
		},
	}

	if strings.Contains(i, "@") {
		u.Email = i
	} else {
		u.ID = i
	}
	return u
}

func (u *user) getID() (string, error) {
	var k []byte
	var err error
	if u.ID != "" {
		k, err = u.DS.dcbIndUsers.getBytes(u.ID)
		if err != nil {
			return "", err
		}
	}
	if len(k) == 0 && u.Email != "" {
		k, err = u.DS.dcbIndUsers.getBytes(u.Email)
		if err != nil {
			return "", err
		}
	}
	if len(k) == 0 {
		return "", errors.New("not found")
	}
	return string(k), nil
}

func (u *user) get() (bool, error) {
	id, err := u.getID()
	if err != nil {
		return false, err
	}

	b, err := u.DS.dcbUsers.getBytes(id)
	if len(b) == 0 || err != nil {
		return false, err
	}
	br := bytes.NewReader(b)
	dec := gob.NewDecoder(br)
	if err := dec.Decode(u); err != nil {
		return false, err
	}
	return true, nil
}

func (u *user) set() error {
	if u.Email == "" {
		return errors.New("user data incomplete: missing email")
	}
	id, err := u.getID()
	if err != nil {
		return err
	}

	if err = u.DS.dcbIndUsers.setBytes(u.ID, []byte(id)); err != nil {
		return err
	}
	if err = u.DS.dcbIndUsers.setBytes(u.Email, []byte(id)); err != nil {
		return err
	}

	bb := &bytes.Buffer{}
	enc := gob.NewEncoder(bb)
	if err := enc.Encode(u); err != nil {
		return err
	}
	if err = u.DS.dcbUsers.set(id, bb); err != nil {
		return err
	}
	return nil
}

type post struct {
	*dtPost
}

func (n *node) newPost() *post {
	return &post{
		dtPost: newPost(),
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
			ID: i, DS: n.su.ds,
		},
	}
}

func (ps *posts) get() (bool, error) {
	id, err := ps.getID()
	if err != nil {
		return false, err
	}
	b, err := ps.DS.dcbPosts.getBytes(id)
	if len(b) == 0 || err != nil {
		return false, err
	}
	br := bytes.NewReader(b)
	dec := gob.NewDecoder(br)
	if err := dec.Decode(ps); err != nil {
		return false, err
	}
	return true, nil
}

func (ps *posts) set() error {
	v, err := json.Marshal(ps)
	if err != nil {
		return err
	}
	id, err := ps.getID()
	if err != nil {
		return err
	}
	if err = ps.DS.dcbPosts.setBytes(id, v); err != nil {
		return err
	}
	return nil
}
