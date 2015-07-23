package main

import (
	"bytes"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

func (n *node) postHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	/*
		unreg -

		reg -
			store submitted form for later sending (encrypt)
			get email and send form
	*/
	u := n.newUser()
	u.ID = "test"
	u.PublicID = u.ID
	u.Email = "test@test.com"
	u.Confirm = &dtConfirm{}
	if err := u.set(); err != nil {
		fmt.Println(err)
	}
	rdr, err := n.su.ds.dcbAsts.get("test")
	if err != nil {
		n.su.logs.Err.Println(err)
	}
	b := &bytes.Buffer{}
	b.ReadFrom(rdr)
	fmt.Fprintf(w, string(b.Bytes()))
}

func (n *node) nilHandler(w http.ResponseWriter, r *http.Request) {}
