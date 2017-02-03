package api

import (
	"context"
	"net/http"

	"github.com/codemodus/formlark/internal/entities"
)

// UserProvider ...
type UserProvider interface {
	InsUserClaim(context.Context, *entities.UserRequiz) (*entities.Empty, error)
	SrchUser(context.Context, *entities.UserReferral) (*entities.User, error)
}

func (a *API) userClaimPostHandler(w http.ResponseWriter, r *http.Request) {
	ur := &entities.UserRequiz{}

	if err := decodeBody(r, ur); err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}

	e, err := a.userP.InsUserClaim(r.Context(), ur)
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
	}

	u, err := a.userP.SrchUser(r.Context(), ur)
	if err != nil {
		httpError(w, http.StatusNotFound)
		return
	}

	if err = encodeBody(w, u); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
}
