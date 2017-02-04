package api

import (
	"context"
	"net/http"

	"github.com/codemodus/formlark/internal/entities"
	"github.com/codemodus/formlark/internal/httperr"
)

// UserProvider ...
type UserProvider interface {
	InsUserClaim(context.Context, *entities.UserRequiz) (*entities.Empty, httperr.HTTPError)
	SrchUser(context.Context, *entities.UserReferral) (*entities.User, httperr.HTTPError)
}

func (a *API) userClaimPostHandler(w http.ResponseWriter, r *http.Request) {
	ur := &entities.UserRequiz{}

	if err := decodeBody(r, ur); err != nil {
		errorHandler(w, httperr.New(err, http.StatusBadRequest, ""))
		return
	}

	e, herr := a.userP.InsUserClaim(r.Context(), ur)
	if herr != nil {
		errorHandler(w, herr)
		return
	}

	if err := encodeBody(w, e); err != nil {
		errorHandler(w, httperr.New(err, http.StatusInternalServerError, ""))
		return
	}
}

func (a *API) userGetSearchHandler(w http.ResponseWriter, r *http.Request) {
	ur := &entities.UserReferral{
		Email: r.URL.Query().Get("email"),
	}

	u, herr := a.userP.SrchUser(r.Context(), ur)
	if herr != nil {
		errorHandler(w, herr)
		return
	}

	if err := encodeBody(w, u); err != nil {
		errorHandler(w, httperr.New(err, http.StatusInternalServerError, ""))
		return
	}
}
