package main

import "time"

type dtConfirm struct {
	Email string            `json:"email,omitempty"`
	Forms map[string]string `json:"forms,omitempty"`
}

type dtUser struct {
	Email   string     `json:"email"`
	Forms   []string   `json:"forms,omitempty"`
	Confirm *dtConfirm `json:"confirm,omitempty"`
}

type dtPost struct {
	Date    time.Time           `json:"date"`
	Replyto string              `json:"replyto"`
	Next    string              `json:"next"`
	Subject string              `json:"subject"`
	CC      []string            `json:"cc"`
	Content map[string][]string `json:"content"`
}
