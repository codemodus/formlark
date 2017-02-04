package inmem

import (
	"fmt"
	"time"

	"github.com/codemodus/formlark/internal/entities"
	"github.com/codemodus/uidg"
)

type claimedUser struct {
	Token        uint64
	ExpirationAt time.Time
	User         *entities.User
}

// InMem ...
type InMem struct {
	idg          *uidg.UIDG
	auths        map[uint64]struct{}
	users        map[uint64]*entities.User
	claimedUsers map[uint64]*claimedUser
}

// New ...
func New() (*InMem, error) {
	im := &InMem{
		idg:          uidg.New(),
		auths:        make(map[uint64]struct{}),
		users:        make(map[uint64]*entities.User),
		claimedUsers: make(map[uint64]*claimedUser),
	}

	testAuth := im.idg.Gen()
	fmt.Printf("ADDING TEST AUTHORIZATION TOKEN: %d\n", testAuth)
	im.auths[testAuth] = struct{}{}

	return im, nil
}
