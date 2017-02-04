package inmem

import (
	"context"
	"net/http"
	"time"

	"github.com/codemodus/formlark/internal/entities"
	"github.com/codemodus/formlark/internal/httperr"
)

// InsMessageByUserID ...
func (i *InMem) InsMessageByUserID(ctx context.Context, mr *entities.MessageByUserIDRecord) (*entities.Message, httperr.HTTPError) {
	u, ok := i.users[mr.UserID]
	if !ok {
		return nil, httperr.New(nil, http.StatusNotFound, "cannot send msg to non-existent user")
	}

	if u.ConfirmedAt.Time.After(time.Now()) {
		return nil, httperr.New(nil, http.StatusUnauthorized, "user must first be confirmed")
	}

	m := &entities.Message{
		ID:        i.idg.Gen(),
		UserID:    mr.UserID,
		CreatedAt: time.Now(),
		Form:      make(map[string]string),
	}

	for k, v := range mr.Message.Form {
		if k == "_replyto" {
			m.ReplyTo = v
			continue
		}
		if k == "_subject" {
			m.Subject = v
			continue
		}
		m.Form[k] = v
	}

	return m, nil
}
