package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codemodus/formlark/internal/cx"
)

func (a *API) reco(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// TODO: add logging
				fmt.Println(err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *API) authCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			a, err := strconv.ParseUint(auth, 10, 64)
			if err != nil {
				httpError(w, http.StatusBadRequest)
			}

			cx.ReqSetHTTPAuth(r, a)
		}

		tAuth := r.Header.Get("Temp-Authorization")
		if tAuth != "" {
			ta, err := strconv.ParseUint(tAuth, 10, 64)
			if err != nil {
				httpError(w, http.StatusBadRequest)
			}

			cx.ReqSetHTTPTempAuth(r, ta)
		}

		next.ServeHTTP(w, r)
	})
}
