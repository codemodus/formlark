package sessmgr

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Provider interface {
	Create(string) (*Session, error)
	Read(string) (*Session, error)
	Update(string) error
	Destroy(string) error
	GC(int64)
	// TODO: GC and persistent providers.
}

type VolatileProvider struct {
	mu       sync.RWMutex
	sessions map[string]*list.Element
	list     *list.List
}

func NewVolatileProvider() *VolatileProvider {
	s := make(map[string]*list.Element)
	return &VolatileProvider{
		sessions: s,
		list:     list.New(),
	}
}

// TODO: typed errors
func (p *VolatileProvider) Create(id string) (*Session, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	s := NewSession(id, p)

	elem := p.list.PushBack(s)
	p.sessions[id] = elem

	return s, nil
}

func (p *VolatileProvider) Read(id string) (*Session, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if elem, ok := p.sessions[id]; ok {
		return elem.Value.(*Session), nil
	}

	return nil, fmt.Errorf(`session %q does not exist`, id)
}

func (p *VolatileProvider) Update(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if elem, ok := p.sessions[id]; ok {
		elem.Value.(*Session).Last = time.Now()
		p.list.MoveToFront(elem)
	}

	return nil
}

func (p *VolatileProvider) Destroy(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if elem, ok := p.sessions[id]; ok {
		delete(p.sessions, id)
		p.list.Remove(elem)
	}

	return nil
}

func (p *VolatileProvider) GC(maxLife int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for elem := p.list.Back(); elem != nil; elem = p.list.Back() {
		if (elem.Value.(*Session).Last.Unix() + maxLife) < time.Now().Unix() {
			p.list.Remove(elem)
			delete(p.sessions, elem.Value.(*Session).ID())
		} else {
			continue
		}
	}
}
