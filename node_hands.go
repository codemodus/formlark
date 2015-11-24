package main

import (
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/daved/context"
)

func (n *node) iconHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "front/assets/public/icon/"+r.URL.Path)
}

func (n *node) assetsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:len(n.u.conf.AdminPathPrefix)+1] == n.u.conf.AdminPathPrefix {
		http.ServeFile(w, r, "front/assets/protected/"+r.URL.Path[9+len(n.u.conf.AdminPathPrefix):])
		return
	}
	http.ServeFile(w, r, "front/assets/public/"+r.URL.Path[8:])
}

func (n *node) assetsFlexHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:len(n.u.conf.AdminPathPrefix)+1] == n.u.conf.AdminPathPrefix {
		http.ServeFile(w, r, "front/assets/"+r.URL.Path[1+len(n.u.conf.AdminPathPrefix):])
		return
	}
	http.ServeFile(w, r, "front/assets/"+r.URL.Path[1:])
}

func (n *node) backupHandleFunc(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	err := n.u.ds.dcbsRsrcs.DB.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (n *node) NotFound(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
