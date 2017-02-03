package api

import (
	"net/http"

	"github.com/codemodus/formlark/internal/entities"
	"github.com/codemodus/parth"
)

// MessageProvider ...
type MessageProvider interface {
	InsMessageByUserID(*entities.MessageByUserIDRecord) (*entities.Message, error)
}

func (a *API) messageByUserPostHandler(w http.ResponseWriter, r *http.Request) {
	mr := &entities.MessageByUserIDRecord{}

	if err := decodeBody(r, mr); err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}

	uid, err := parth.SubSegToInt64(r.URL.Path, "user")
	if err != nil {
		httpError(w, http.StatusNotFound)
		return
	}
	mr.UserID = uint64(uid)

	m, err := a.msgP.InsMessageByUserID(mr)
	if err != nil {
		panic(err)
	}

	if err = encodeBody(w, m); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
}
