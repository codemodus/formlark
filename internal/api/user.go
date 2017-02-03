package api

import (
	"net/http"

	"github.com/codemodus/formlark/internal/entities"
)

// UserProvider ...
type UserProvider interface {
	InsUserClaim(*entities.UserRequiz) (*entities.Empty, error)
	SrchUser(*entities.UserReferral) (*entities.User, error)
}

func (a *API) userClaimPostHandler(w http.ResponseWriter, r *http.Request) {
	ur := &entities.UserRequiz{}

	if err := decodeBody(r, ur); err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}

	e, err := a.userP.InsUserClaim(ur)
	if err != nil {
		panic(err)
	}

	if err = encodeBody(w, e); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
}

func (a *API) userGetSearchHandler(w http.ResponseWriter, r *http.Request) {
	ur := &entities.UserReferral{
		Email: r.URL.Query().Get("email"),
		Token: r.URL.Query().Get("token"),
	}

	u, err := a.userP.SrchUser(ur)
	if err != nil {
		panic(err)
	}

	if err = encodeBody(w, u); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
}
