package inmem

import (
	"context"
	"fmt"
	"time"

	"github.com/codemodus/formlark/internal/entities"
)

// InsMessageByUserID ...
func (i *InMem) InsMessageByUserID(ctx context.Context, mr *entities.MessageByUserIDRecord) (*entities.Message, error) {
	u, ok := i.users[mr.UserID]
	if !ok {
		return nil, fmt.Errorf("no user by id to rcv msg")
	}

	if u.ConfirmedAt.Time.After(time.Now()) {
		return nil, fmt.Errorf("user is not confirmed yet")
	}

	m := &entities.Message{
		ID:        i.idg.Gen(),
		UserID:    mr.UserID,
		CreatedAt: time.Now(),
	}

	for k, v := range mr.Message.Form {
		if k == "_replyto" {
			m.ReplyTo = v
			delete(mr.Message.Form, k)
		}
		if k == "_subject" {
			m.Subject = v
			delete(mr.Message.Form, k)
		}
	}

	m.Form = mr.Message.Form

	return m, nil
}
