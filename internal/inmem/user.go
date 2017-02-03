package inmem

import (
	"context"
	"fmt"
	"time"

	"github.com/codemodus/formlark/internal/cx"
	"github.com/codemodus/formlark/internal/entities"
)

// InsUserClaim ...
func (i *InMem) InsUserClaim(ctx context.Context, ur *entities.UserRequiz) (*entities.Empty, error) {
	for _, v := range i.users {
		if v.Email == ur.User.Email {
			return nil, fmt.Errorf("user exists")
		}
	}

	t := time.Now()

	cu := &claimedUser{
		Token:        i.idg.Gen(),
		ExpirationAt: time.Now().Add(time.Minute * 2),
		User: &entities.User{
			ID:        i.idg.Gen(),
			Email:     ur.User.Email,
			CreatedAt: t,
			UpdatedAt: t,
		},
	}

	fmt.Println(cu.Token, cu.User.Email)

	i.claimedUsers[cu.User.ID] = cu

	return &entities.Empty{}, nil
}

// SrchUser ...
func (i *InMem) SrchUser(ctx context.Context, ur *entities.UserReferral) (*entities.User, error) {
	var u *entities.User

	t, ok := cx.HTTPTempAuth(ctx)
	if ok {
		return i.srchUserClaim(t, ur)
	}

	a, ok := cx.HTTPAuth(ctx)
	if !ok || !i.isValidAuth(a) {
		return nil, fmt.Errorf("not authorized")
	}

	for _, v := range i.users {
		if v.Email == ur.Email {
			u = i.users[v.ID]
		}
	}

	if u == nil {
		return nil, fmt.Errorf("no user found")
	}

	return u, nil
}

func (i *InMem) srchUserClaim(token uint64, ur *entities.UserReferral) (*entities.User, error) {
	var cu *claimedUser

	for _, v := range i.claimedUsers {
		if v.User.Email == ur.Email {
			cu = i.claimedUsers[v.User.ID]
		}
	}

	if cu == nil {
		return nil, nil
	}

	if token == 0 || cu.Token != token {
		return nil, fmt.Errorf("bad token")
	}

	if cu.ExpirationAt.Before(time.Now()) {
		return nil, fmt.Errorf("expired token")
	}

	cu.User.ConfirmedAt.Time = time.Now()
	cu.User.ConfirmedAt.Valid = true

	i.users[cu.User.ID] = cu.User

	delete(i.claimedUsers, cu.User.ID)

	return cu.User, nil
}
