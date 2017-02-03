package inmem

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codemodus/formlark/internal/entities"
)

// InsUserClaim ...
func (i *InMem) InsUserClaim(ur *entities.UserRequiz) (*entities.Empty, error) {
	for _, v := range i.users {
		if v.Email == ur.User.Email {
			return nil, fmt.Errorf("user exists")
		}
	}

	t := time.Now()

	u := &entities.User{
		ID:        i.idg.Gen(),
		Email:     ur.User.Email,
		CreatedAt: t,
		UpdatedAt: t,
		Token:     strconv.FormatUint(i.idg.Gen(), 10),
	}

	i.users[u.ID] = u

	return &entities.Empty{}, nil
}

// SrchUser ...
func (i *InMem) SrchUser(ur *entities.UserReferral) (*entities.User, error) {
	var u *entities.User

	for _, v := range i.users {
		if v.Email == ur.Email {
			u = i.users[v.ID]
		}
	}

	if u == nil {
		return nil, fmt.Errorf("no user found")
	}

	if u.Token != ur.Token {
		return nil, fmt.Errorf("bad token")
	}

	if ur.Token != "" && ur.Token == u.Token {
		u.Token = ""
		u.ConfirmedAt = entities.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	return u, nil
}
