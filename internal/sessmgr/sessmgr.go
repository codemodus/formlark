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
	SessCreate(string) (Session, error)
	SessRead(string) (Session, error)
	SessUpdate(string) error
	SessDestroy(string)
	SessGC(int64)
}

type ProviderRegistry map[string]Provider

func NewProviderRegistry() *ProviderRegistry {
	m := make(map[string]Provider)
	p := ProviderRegistry(m)
	return &p
}

func (p *ProviderRegistry) Register(name string, pdr Provider) {
	if pdr == nil {
		panic("session manager: register called with provider as nil")
	}
	if _, dup := (*p)[name]; dup {
		panic(`session manager: register called twice for provider "` + name + `"`)
	}
	(*p)[name] = pdr
}

type Manager struct {
	Mu      sync.Mutex
	ProReg  *ProviderRegistry
	Name    string
	Pvd     Provider
	MaxLife int64
}

func New(proReg *ProviderRegistry, providerName, cookieName string, maxLife int64) (*Manager, error) {
	pvd, ok := (*proReg)[providerName]
	if !ok {
		return nil, fmt.Errorf("session manager: unknown provider %q", providerName)
	}
	return &Manager{ProReg: proReg, Pvd: pvd, Name: cookieName, MaxLife: maxLife}, nil
}

func (sm *Manager) SessID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (sm *Manager) SessStart(w http.ResponseWriter, r *http.Request) (s Session) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	c, err := r.Cookie(sm.Name)
	if err == nil && c.Value != "" {
		id, _ := url.QueryUnescape(c.Value)
		s, err = sm.Pvd.SessRead(id)
		if err == nil {
			return s
		}
	}

	id := sm.SessID()
	if s, err = sm.Pvd.SessCreate(id); err != nil {
		panic("how could you do this?")
	}
	c = &http.Cookie{Name: sm.Name, Value: url.QueryEscape(id), Path: "/", HttpOnly: true, MaxAge: int(sm.MaxLife)}
	http.SetCookie(w, c)
	return s
}

type Sess struct {
	ID      string
	LastAcs time.Time
	Val     map[string]interface{}
	Pvd     Provider
}

func (s *Sess) Set(key string, val interface{}) error {
	s.Val[key] = val
	s.Pvd.SessUpdate(s.ID)
	return nil
}

func (s *Sess) Get(key string) interface{} {
	s.Pvd.SessUpdate(s.ID)
	if v, ok := s.Val[key]; ok {
		return v
	}
	return nil
}

func (s *Sess) Delete(key string) {
	delete(s.Val, key)
	s.Pvd.SessUpdate(s.ID)
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

func (p *VolatileProvider) SessCreate(id string) (Session, error) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	v := make(map[string]interface{}, 0)
	s := &Sess{ID: id, LastAcs: time.Now(), Val: v, Pvd: p}
	elem := p.list.PushBack(s)
	p.sessions[id] = elem
	return s, nil
}

func (p *VolatileProvider) SessRead(id string) (Session, error) {
	if elem, ok := p.sessions[id]; ok {
		return elem.Value.(*Sess), nil
	}
	return nil, fmt.Errorf(`session %q does not exist`, id)
}

func (p *VolatileProvider) SessDestroy(id string) {
	if elem, ok := p.sessions[id]; ok {
		delete(p.sessions, id)
		p.list.Remove(elem)
	}
}

func (p *VolatileProvider) SessGC(maxLife int64) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for elem := p.list.Back(); elem != nil; elem = p.list.Back() {
		if (elem.Value.(*Sess).LastAcs.Unix() + maxLife) < time.Now().Unix() {
			p.list.Remove(elem)
			delete(p.sessions, elem.Value.(*Sess).ID)
		} else {
			continue
		}
	}
}

func (p *VolatileProvider) SessUpdate(id string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if elem, ok := p.sessions[id]; ok {
		elem.Value.(*Sess).LastAcs = time.Now()
		p.list.MoveToFront(elem)
	}
	return nil
}
