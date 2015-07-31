package main

import (
	"errors"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

type dtConfirm struct {
	Forms map[string]string `json:"forms,omitempty"`
}

func newConfirm() *dtConfirm {
	return &dtConfirm{Forms: make(map[string]string)}
}

type dtUser struct {
	Email   string     `json:"email"`
	Forms   []string   `json:"forms,omitempty"`
	Confirm *dtConfirm `json:"confirm,omitempty"`
}

func newUser() *dtUser {
	return &dtUser{
		Forms: make([]string, 0), Confirm: newConfirm(),
	}
}

func (u *dtUser) validate() error {
	if u.Email == "" {
		return errors.New("email must not be blank")
	}
	t := ""
	if !strings.Contains(u.Email, "<") || !strings.Contains(u.Email, ">") {
		t = u.Email + " <" + u.Email + ">"
	}
	a, err := mail.ParseAddress(t)
	if err != nil {
		return err
	}
	u.Email = a.Address

	return nil
}

type dtPost struct {
	Date    time.Time           `json:"date"`
	Referer *url.URL            `json:"url"`
	Replyto string              `json:"replyto"`
	Next    string              `json:"next"`
	Subject string              `json:"subject"`
	CC      []string            `json:"cc"`
	Content map[string][]string `json:"content"`
}

func newPost() *dtPost {
	return &dtPost{CC: make([]string, 0), Content: make(map[string][]string)}
}

func (p *dtPost) validate() error {
	if p.Replyto != "" {
		t := ""
		if !strings.Contains(p.Replyto, "<") || !strings.Contains(p.Replyto, ">") {
			t = p.Replyto + " <" + p.Replyto + ">"
		}
		a, err := mail.ParseAddress(t)
		if err != nil {
			return err
		}
		p.Replyto = a.Address
	}
	if len(p.CC) > 0 {
		for k, v := range p.CC {
			t := ""
			if !strings.Contains(v, "<") || !strings.Contains(v, ">") {
				t = v + " <" + v + ">"
			}
			a, err := mail.ParseAddress(t)
			if err != nil {
				return err
			}
			p.CC[k] = a.Address
		}
	}

	return nil
}
