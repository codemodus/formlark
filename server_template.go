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
	ItemsTitle   string
	ItemsCommon  []NavItem
	ItemsSpecial []NavItem
}

type Footer struct {
	ColsDropdownFlag bool
	ColsDropdown     []NavGroup
	RowBottom        NavGroup
}

type Page struct {
	AppName   string
	PageTitle string
	URLLogin  string
	NavHeader NavGroup
	NavDrawer NavGroup
	Footer    Footer
	Misc      string
}

func newPage() *Page {
	p := &Page{AppName: "Formlark", PageTitle: "Formlark", URLLogin: "/login"}
	return p
}

func (n *node) newNavFooterBottomItemsCommon() []NavItem {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavItem, 3, 3)
	r[0] = NavItem{HRef: "/support", Name: "Support"}
	r[1] = NavItem{HRef: "/privacy", Name: "Privacy"}
	r[2] = NavItem{HRef: "/legal", Name: "Legal"}
	return r
}

type PageAnon struct {
	*Page
}

func (n *node) newPageAnon() *PageAnon {
	p := newPage()
	p.NavHeader.ItemsCommon = n.newNavAnonItemsCommon()
	p.NavDrawer.ItemsCommon = n.newNavAnonItemsCommon()
	p.Footer.RowBottom.ItemsCommon = n.newNavFooterBottomItemsCommon()
	return &PageAnon{p}
}

func (n *node) newNavAnonItemsCommon() []NavItem {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavItem, 1, 1)
	r[0] = NavItem{HRef: "/login", Name: "Login"}
	return r
}

type PageAuthed struct {
	*Page
}

func (n *node) newPageAuthed() *PageAuthed {
	p := newPage()
	p.NavHeader.ItemsCommon = n.newNavAuthedItemsCommon()
	p.NavDrawer.ItemsCommon = n.newNavAuthedItemsCommon()
	p.Footer.ColsDropdownFlag = true
	p.Footer.ColsDropdown = n.newNavFooterAuthedColsDdItemsCommon()
	p.Footer.RowBottom.ItemsCommon = n.newNavFooterBottomItemsCommon()
	return &PageAuthed{p}
}

func (n *node) newNavAuthedItemsCommon() []NavItem {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavItem, 3, 3)
	r[0] = NavItem{HRef: "/overview", Name: "Overview"}
	r[1] = NavItem{HRef: "/settings", Name: "Settings"}
	r[2] = NavItem{HRef: "/logout", Name: "Logout"}
	return r
}

func (n *node) newNavFooterAuthedColsDdItemsCommon() []NavGroup {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavGroup, 3, 3)
	r[0] = NavGroup{ItemsTitle: "Test", ItemsCommon: n.newNavFooterBottomItemsCommon()}
	r[1] = NavGroup{ItemsTitle: "Rest", ItemsCommon: n.newNavFooterBottomItemsCommon()}
	r[2] = NavGroup{ItemsTitle: "Best", ItemsCommon: n.newNavFooterBottomItemsCommon()}
	return r
}

type PageAdmin struct {
	*Page
}

func (n *node) newPageAdmin() *PageAdmin {
	p := newPage()
	p.URLLogin = "/" + n.su.conf.AdminPathPrefix + "/login"
	p.NavHeader.ItemsCommon = n.newNavAdminItemsCommon()
	p.NavDrawer.ItemsCommon = n.newNavAdminItemsCommon()
	p.Footer.ColsDropdownFlag = true
	p.Footer.ColsDropdown = n.newNavFooterAuthedColsDdItemsCommon()
	p.Footer.RowBottom.ItemsCommon = n.newNavFooterBottomItemsCommon()
	return &PageAdmin{p}
}

func (n *node) newNavAdminItemsCommon() []NavItem {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavItem, 4, 4)
	r[0] = NavItem{HRef: "/" + n.su.conf.AdminPathPrefix + "/overview", Name: "Overview"}
	r[1] = NavItem{HRef: "/" + n.su.conf.AdminPathPrefix + "/users", Name: "Users"}
	r[2] = NavItem{HRef: "/" + n.su.conf.AdminPathPrefix + "/settings", Name: "Settings"}
	r[3] = NavItem{HRef: "/" + n.su.conf.AdminPathPrefix + "/logout", Name: "Logout"}
	return r
}

func (n *node) newNavFooterAdminColsDdItemsCommon() []NavGroup {
	// TODO: Move common items to node field and init within node setup.
	r := make([]NavGroup, 3, 3)
	r[0] = NavGroup{ItemsTitle: "Test", ItemsCommon: n.newNavFooterBottomItemsCommon()}
	r[1] = NavGroup{ItemsTitle: "Rest", ItemsCommon: n.newNavFooterBottomItemsCommon()}
	r[2] = NavGroup{ItemsTitle: "Best", ItemsCommon: n.newNavFooterBottomItemsCommon()}
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
