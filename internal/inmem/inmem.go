package inmem

import (
	"github.com/codemodus/formlark/internal/entities"
	"github.com/codemodus/uidg"
)

// InMem ...
type InMem struct {
	idg   *uidg.UIDG
	users map[uint64]*entities.User
}

// New ...
func New() (*InMem, error) {
	im := &InMem{
		idg:   uidg.New(),
		users: make(map[uint64]*entities.User),
	}

	return im, nil
}
