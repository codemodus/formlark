package sessmgr

import (
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func unsetCookie(w http.ResponseWriter, name string) {
	c := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, c)
}

type Session interface {
	Set(string, interface{}) error
	Get(string) interface{}
	Unset(string)
	SessID() string
}

type Provider interface {
	Create(string) (Session, error)
	Read(string) (Session, error)
	Update(string) error
	Destroy(string) error
	GC(int64)
	// TODO: GC and persistent providers.
}

type Manager struct {
	name    string
	maxLife int64
	prov    Provider
}

func New(name string, maxLife int64, provider Provider) *Manager {
	return &Manager{
		name:    name,
		maxLife: maxLife,
		prov:    provider,
	}
}

func (m *Manager) genSessID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (m *Manager) SessStart(w http.ResponseWriter, r *http.Request) (s Session, err error) {
	c, err := r.Cookie(m.name)
	if err == nil {
		id, err := url.QueryUnescape(c.Value)
		if err == nil && id != "" {
			s, err = m.prov.Read(id)
			if err == nil {
				return s, nil
			}
		}
		unsetCookie(w, m.name)
	}

	id := m.genSessID()
	if s, err = m.prov.Create(id); err != nil {
		return nil, err
	}
	c = &http.Cookie{Name: m.name, Value: url.QueryEscape(id), Path: "/", HttpOnly: true, MaxAge: int(m.maxLife)}
	http.SetCookie(w, c)
	return s, nil
}

func (m *Manager) SessStop(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(m.name)
	if err != nil {
		return
	}
	unsetCookie(w, m.name)
	id, err := url.QueryUnescape(c.Value)
	if err != nil || id != "" {
		return
	}
	m.prov.Destroy(id)
}

type Sess struct {
	ID   string
	Last time.Time
	mu   sync.RWMutex
	Val  map[string]interface{}
	prov Provider
}

func NewSess(id string, provider Provider) *Sess {
	v := make(map[string]interface{}, 0)
	return &Sess{ID: id, Last: time.Now(), Val: v, prov: provider}
}

func (s *Sess) Set(key string, val interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Val[key] = val
	s.prov.Update(s.ID)
	return nil
}

func (s *Sess) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.prov.Update(s.ID)
	v, _ := s.Val[key]
	return v
}

func (s *Sess) Unset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Val, key)
	s.prov.Update(s.ID)
}

func (s *Sess) SessID() string {
	return s.ID
}

type VolatileProvider struct {
	mu       sync.RWMutex
	sessions map[string]*list.Element
	list     *list.List
}

func NewVolatileProvider() *VolatileProvider {
	s := make(map[string]*list.Element)
	return &VolatileProvider{sessions: s, list: list.New()}
}

// TODO: typed errors
func (p *VolatileProvider) Create(id string) (Session, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	s := NewSess(id, p)
	elem := p.list.PushBack(s)
	p.sessions[id] = elem
	return s, nil
}

func (p *VolatileProvider) Read(id string) (Session, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if elem, ok := p.sessions[id]; ok {
		return elem.Value.(*Sess), nil
	}
	return nil, fmt.Errorf(`session %q does not exist`, id)
}

func (p *VolatileProvider) Update(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if elem, ok := p.sessions[id]; ok {
		elem.Value.(*Sess).Last = time.Now()
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
		if (elem.Value.(*Sess).Last.Unix() + maxLife) < time.Now().Unix() {
			p.list.Remove(elem)
			delete(p.sessions, elem.Value.(*Sess).ID)
		} else {
			continue
		}
	}
}
