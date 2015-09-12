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

type NavGroup struct {
	NavCommonItems  []NavItem
	NavSpecialItems []NavItem
}

type Page struct {
	AppName   string
	URLLogin  string
	NavHeader NavGroup
	NavDrawer NavGroup
	Misc      string
}

type PagePublic struct {
	*Page
}

func newPage() *Page {
	return &Page{AppName: "Formlark", URLLogin: "/login"}
}

func (n *node) newPagePublic() *PagePublic {
	p := newPage()
	p.NavHeader.NavCommonItems = n.newNavCommonItems()
	p.NavDrawer.NavCommonItems = n.newNavCommonItems()
	return &PagePublic{p}
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

type PageAdmin struct {
	*Page
}

func (n *node) newPageAdmin() *PageAdmin {
	p := newPage()
	p.URLLogin = "/" + n.su.conf.AdminPathPrefix + "/login"
	p.NavHeader.NavCommonItems = n.newNavAdminCommonItems()
	p.NavDrawer.NavCommonItems = n.newNavAdminCommonItems()
	return &PageAdmin{p}
}

func (n *node) newNavAdminCommonItems() []NavItem {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavItem, 5, 5)
	r[0] = NavItem{HRef: "/atada", Name: "aTada!"}
	r[1] = NavItem{HRef: "/alada", Name: "aLada!"}
	r[2] = NavItem{HRef: "/amada", Name: "aMada!"}
	r[3] = NavItem{HRef: "/arada", Name: "aRada!"}
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
