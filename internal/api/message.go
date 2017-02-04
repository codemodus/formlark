package api

import (
	"context"
	"net/http"

	"github.com/codemodus/formlark/internal/entities"
	"github.com/codemodus/formlark/internal/httperr"
	"github.com/codemodus/parth"
)

// MessageProvider ...
type MessageProvider interface {
	InsMessageByUserID(context.Context, *entities.MessageByUserIDRecord) (*entities.Message, httperr.HTTPError)
}

func (a *API) messageByUserPostHandler(w http.ResponseWriter, r *http.Request) {
	mr := &entities.MessageByUserIDRecord{}

	if err := decodeBody(r, mr); err != nil {
		errorHandler(w, httperr.New(err, http.StatusBadRequest, ""))
		return
	}

	uid, err := parth.SubSegToInt64(r.URL.Path, "user")
	if err != nil {
		errorHandler(w, httperr.New(err, http.StatusNotFound, ""))
		return
	}
	mr.UserID = uint64(uid)

	m, herr := a.msgP.InsMessageByUserID(r.Context(), mr)
	if herr != nil {
		errorHandler(w, herr)
		return
	}

	if err = encodeBody(w, m); err != nil {
		errorHandler(w, httperr.New(err, http.StatusInternalServerError, ""))
		return
	}
}
