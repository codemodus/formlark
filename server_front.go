package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

func (n *node) postHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	/*
		unreg -

		reg -
			store submitted form for later sending (encrypt)
			get email and send form
	*/
	n.su.ds.dcbAsts.setBytes("test", []byte(strconv.FormatInt(time.Now().Unix(), 10)))
	rdr, err := n.su.ds.dcbAsts.get("test")
	if err != nil {
		n.su.logs.Err.Println(err)
	}
	b := &bytes.Buffer{}
	b.ReadFrom(rdr)
	fmt.Fprintf(w, string(b.Bytes()))
}

func (n *node) nilHandler(w http.ResponseWriter, r *http.Request) {}
