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
	PageTitle string
	URLLogin  string
	NavHeader NavGroup
	NavDrawer NavGroup
	Misc      string
}

type PagePublic struct {
	*Page
}

func newPage() *Page {
	return &Page{AppName: "Formlark", PageTitle: "Formlark", URLLogin: "/login"}
}

func (n *node) newPagePublic() *PagePublic {
	p := newPage()
	p.NavHeader.NavCommonItems = n.newNavCommonItems()
	p.NavDrawer.NavCommonItems = n.newNavCommonItems()
	return &PagePublic{p}
}

func (n *node) newNavCommonItems() []NavItem {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavItem, 3, 3)
	r[0] = NavItem{HRef: "/overview", Name: "Overview"}
	r[1] = NavItem{HRef: "/settings", Name: "Settings"}
	r[2] = NavItem{HRef: "/logout", Name: "Logout"}
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
	r := make([]NavItem, 4, 4)
	r[0] = NavItem{HRef: "/" + n.su.conf.AdminPathPrefix + "/overview", Name: "Overview"}
	r[1] = NavItem{HRef: "/" + n.su.conf.AdminPathPrefix + "/users", Name: "Users"}
	r[2] = NavItem{HRef: "/" + n.su.conf.AdminPathPrefix + "/settings", Name: "Settings"}
	r[3] = NavItem{HRef: "/" + n.su.conf.AdminPathPrefix + "/logout", Name: "Logout"}
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
