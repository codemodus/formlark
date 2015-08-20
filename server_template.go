package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func (n *node) newPage() *Page {
	p := NewPage()
	p.URLLogin = "/" + n.su.conf.AdminPathPrefix + "/login"
	return p
}

func (n *node) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	b := &bytes.Buffer{}
	err := n.su.ts.m[name].ExecuteTemplate(b, "", data)
	if err != nil {
		// TODO: Log
		fmt.Println(err)
		b.Reset()
		b = &bytes.Buffer{}
		http.Error(w, "template failed - please contact the site admin", 500)
		return
	}
	_, err = b.WriteTo(w)
	if err != nil {
		// TODO: Log
		fmt.Println(err)
	}
}
