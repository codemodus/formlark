package api

import (
	"net/http"

	"github.com/codemodus/chain"
	"github.com/codemodus/mixmux"
)

// DataProvider ...
type DataProvider interface {
	UserProvider
	MessageProvider
}

// API ...
type API struct {
	mux   http.Handler
	userP UserProvider
	msgP  MessageProvider
}

// New ...
func New(dp DataProvider) (*API, error) {
	a := &API{
		userP: dp,
		msgP:  dp,
	}

	a.setMux()

	return a, nil
}

func (a *API) setMux() {
	m := mixmux.NewTreeMux(nil)

	c := chain.New(a.reco, a.authCtx)

	m.Get("/user", c.EndFn(a.userGetSearchHandler))
	m.OptionsHeaders("/user")
	m.Post("/claim/user", c.EndFn(a.userClaimPostHandler))
	m.OptionsHeaders("/claim/user")

	m.Post("/user/:id/message", c.EndFn(a.messageByUserPostHandler))
	m.OptionsHeaders("/user/:id/message")

	a.mux = m
}

// ServeHTTP ...
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
