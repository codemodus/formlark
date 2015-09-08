package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type NavItem struct {
	HRef string
	Name string
}

type NavGroup0 struct {
	NavCommonItems  []NavItem
	NavSpecialItems []NavItem
}

type Page struct {
	AppName   string
	URLLogin  string
	NavHeader NavGroup0
	NavDrawer NavGroup0
	Misc      string
}

func NewPage() *Page {
	return &Page{AppName: "Formlark", URLLogin: "/admin/login"}
}

func (n *node) newPage() *Page {
	p := NewPage()
	p.URLLogin = "/" + n.su.conf.AdminPathPrefix + "/login"
	p.NavHeader.NavCommonItems = n.newNavCommonItems()
	p.NavDrawer.NavCommonItems = n.newNavCommonItems()
	return p
}

func (n *node) newNavCommonItems() []NavItem {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavItem, 5, 5)
	r[0] = NavItem{HRef: "/tada", Name: "Tada!"}
	r[1] = NavItem{HRef: "/lada", Name: "Lada!"}
	r[2] = NavItem{HRef: "/mada", Name: "Mada!"}
	r[3] = NavItem{HRef: "/rada", Name: "Rada!"}
	r[4] = NavItem{HRef: "/zada", Name: "Zada!"}
	return r
}

func (n *node) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	b := &bytes.Buffer{}
	err := n.su.ts.ExecuteTemplate(b, name, data)
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
