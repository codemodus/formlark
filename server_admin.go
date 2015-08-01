package main

import (
	"net/http"

	"golang.org/x/net/context"
)

func (n *node) adminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := n.su.ts.ExecuteTemplate(w, "admin/login.html", nil)
	if err != nil {
		http.Error(w, "template failed - please contact the site admin", 500)
		return
	}
}
