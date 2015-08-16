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

type Session interface {
	Set(string, interface{}) error
	Get(string) interface{}
	Delete(string)
	SessID() string
}

type Provider interface {
	Create(string) (Session, error)
	Read(string) (Session, error)
	Update(string) error
	Destroy(string)
	GC(int64)
	// TODO: GC and persistent providers.
}

type Manager struct {
	// TODO: Mutex???
	Mu      sync.Mutex
	Name    string
	MaxLife int64
	prov    Provider
}

func New(name string, maxLife int64, provider Provider) (*Manager, error) {
	return &Manager{Name: name, MaxLife: maxLife, prov: provider}, nil
}

func (m *Manager) GenSessID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (m *Manager) SessStart(w http.ResponseWriter, r *http.Request) (s Session) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	c, err := r.Cookie(m.Name)
	if err == nil && c.Value != "" {
		id, _ := url.QueryUnescape(c.Value)
		s, err = m.prov.Read(id)
		if err == nil {
			return s
		}
	}

	id := m.GenSessID()
	if s, err = m.prov.Create(id); err != nil {
		panic("how could you do this?")
	}
	c = &http.Cookie{Name: m.Name, Value: url.QueryEscape(id), Path: "/", HttpOnly: true, MaxAge: int(m.MaxLife)}
	http.SetCookie(w, c)
	return s
}

func (m *Manager) SessStop(w http.ResponseWriter, r *http.Request) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	c, err := r.Cookie(m.Name)
	if err == nil && c.Value != "" {
		id, _ := url.QueryUnescape(c.Value)
		m.prov.Destroy(id)
		c = &http.Cookie{Name: m.Name, MaxAge: -1, Expires: time.Unix(1, 0)}
		http.SetCookie(w, c)
	}
}

type Sess struct {
	ID   string
	Last time.Time
	// TODO: Mutex?
	Val  map[string]interface{}
	prov Provider
}

func NewSess(id string, provider Provider) *Sess {
	v := make(map[string]interface{}, 0)
	return &Sess{ID: id, Last: time.Now(), Val: v, prov: provider}
}

func (s *Sess) Set(key string, val interface{}) error {
	s.Val[key] = val
	s.prov.Update(s.ID)
	return nil
}

func (s *Sess) Get(key string) interface{} {
	s.prov.Update(s.ID)
	if v, ok := s.Val[key]; ok {
		return v
	}
	return nil
}

func (s *Sess) Delete(key string) {
	delete(s.Val, key)
	s.prov.Update(s.ID)
}

func (s *Sess) SessID() string {
	return s.ID
}

type VolatileProvider struct {
	Mu       sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func NewVolatileProvider() *VolatileProvider {
	ss := make(map[string]*list.Element)
	return &VolatileProvider{sessions: ss, list: list.New()}
}

func (p *VolatileProvider) Create(id string) (Session, error) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	s := NewSess(id, p)
	elem := p.list.PushBack(s)
	p.sessions[id] = elem
	return s, nil
}

func (p *VolatileProvider) Read(id string) (Session, error) {
	if elem, ok := p.sessions[id]; ok {
		return elem.Value.(*Sess), nil
	}
	return nil, fmt.Errorf(`session %q does not exist`, id)
}

func (p *VolatileProvider) Update(id string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if elem, ok := p.sessions[id]; ok {
		elem.Value.(*Sess).Last = time.Now()
		p.list.MoveToFront(elem)
	}
	return nil
}

func (p *VolatileProvider) Destroy(id string) {
	if elem, ok := p.sessions[id]; ok {
		delete(p.sessions, id)
		p.list.Remove(elem)
	}
}

func (p *VolatileProvider) GC(maxLife int64) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for elem := p.list.Back(); elem != nil; elem = p.list.Back() {
		if (elem.Value.(*Sess).Last.Unix() + maxLife) < time.Now().Unix() {
			p.list.Remove(elem)
			delete(p.sessions, elem.Value.(*Sess).ID)
		} else {
			continue
		}
	}
}
