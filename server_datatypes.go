package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (n *node) newUser() *user {
	return &user{
		dtUser: newUser(),
		boltItem: &boltItem{
			DS: n.su.ds,
		},
	}
}

type user struct {
	*boltItem
	*dtUser
}

func (u *user) getID() (string, error) {
	if u.ID != "" {
		return u.ID, nil
	}
	if u.Email != "" {
		k, err := u.DS.dcbMrks.getBytes(u.Email)
		if err != nil {
			return "", err
		}
		if len(k) > 0 {
			return string(k), nil
		}
	}
	return "", errors.New("no id")
}

func (u *user) get() error {
	id, err := u.getID()
	if err != nil {
		return err
	}
	v, err := u.DS.dcbAsts.getBytes(id)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(v, u); err != nil {
		return err
	}
	return nil
}

func (u *user) set() error {
	v, err := json.Marshal(u)
	if err != nil {
		return err
	}
	id, err := u.getID()
	if err != nil {
		return err
	}
	if u.Email != "" {
		err = u.DS.dcbMrks.setBytes(u.Email, []byte(id))
		if err != nil {
			return errors.New("cannot save markers")
		}
	}
	if err = u.DS.dcbAsts.setBytes(id, v); err != nil {
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

func (n *node) newPosts() *posts {
	return &posts{
		S: make([]*post, 0),
		boltItem: &boltItem{
			DS: n.su.ds,
		},
	}
}

func (ps *posts) get() error {
	id, err := ps.getID()
	if err != nil {
		return err
	}
	v, err := ps.DS.dcbPosts.getBytes(id)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(v, ps); err != nil {
		return err
	}
	return nil
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
