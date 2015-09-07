package datatypes

import (
	"errors"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

type Confirm struct {
	Forms map[string]string `json:"forms,omitempty"`
}

func NewConfirm() *Confirm {
	return &Confirm{Forms: make(map[string]string)}
}

type User struct {
	Email   string   `json:"email"`
	Forms   []string `json:"forms,omitempty"`
	Confirm *Confirm `json:"confirm,omitempty"`
}

func NewUser() *User {
	return &User{
		Forms: make([]string, 0), Confirm: NewConfirm(),
	}
}

func (u *User) Validate() error {
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

type Post struct {
	Date    time.Time           `json:"date"`
	Referer *url.URL            `json:"url"`
	Replyto string              `json:"replyto"`
	Next    string              `json:"next"`
	Subject string              `json:"subject"`
	CC      []string            `json:"cc"`
	Content map[string][]string `json:"content"`
}

func NewPost() *Post {
	return &Post{CC: make([]string, 0), Content: make(map[string][]string)}
}

func (p *Post) Validate() error {
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
