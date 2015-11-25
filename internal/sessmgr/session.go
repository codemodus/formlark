package sessmgr

import (
	"sync"
	"time"
)

type Session struct {
	id   string
	Last time.Time
	mu   sync.RWMutex
	Val  map[string]interface{}
	prov Provider
}

func NewSession(id string, provider Provider) *Session {
	v := make(map[string]interface{})
	return &Session{
		id:   id,
		Last: time.Now(),
		Val:  v, prov: provider,
	}
}

func (s *Session) Set(key string, val interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Val[key] = val
	s.prov.Update(s.id)

	return nil
}

func (s *Session) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// TODO: ? deal with bool
	v, _ := s.Val[key]
	s.prov.Update(s.id)

	return v
}

func (s *Session) Unset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Val, key)
	s.prov.Update(s.id)
}

func (s *Session) ID() string {
	return s.id
}
