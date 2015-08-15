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
	SessionID() string
}

type Provider interface {
	SessionCreate(string) (Session, error)
	SessionRead(string) (Session, error)
	SessionUpdate(string) error
	SessionDestroy(string)
	SessionGC(int64)
}

type Provisor map[string]Provider

func NewProvisor() *Provisor {
	m := make(map[string]Provider)
	p := Provisor(m)
	return &p
}

func (p *Provisor) Register(name string, pdr Provider) {
	if pdr == nil {
		panic("session manager: register called with provider as nil")
	}
	if _, dup := (*p)[name]; dup {
		panic(`session manager: register called twice for provider "` + name + `"`)
	}
	(*p)[name] = pdr
}

type SessionManager struct {
	Mu      sync.Mutex
	Pvr     *Provisor
	Name    string
	Pvd     Provider
	MaxLife int64
}

func NewSessionManager(provisor *Provisor, providerName, cookieName string, maxLife int64) (*SessionManager, error) {
	pvd, ok := (*provisor)[providerName]
	if !ok {
		return nil, fmt.Errorf("session manager: unknown provider %q", providerName)
	}
	return &SessionManager{Pvr: provisor, Pvd: pvd, Name: cookieName, MaxLife: maxLife}, nil
}

func (sm *SessionManager) SessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (sm *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (s Session) {
	sm.Mu.Lock()
	defer sm.Mu.Unlock()

	c, err := r.Cookie(sm.Name)
	if err == nil && c.Value != "" {
		id, _ := url.QueryUnescape(c.Value)
		s, err = sm.Pvd.SessionRead(id)
		if err == nil {
			return s
		}
	}

	id := sm.SessionID()
	if s, err = sm.Pvd.SessionCreate(id); err != nil {
		panic("how could you do this?")
	}
	c = &http.Cookie{Name: sm.Name, Value: url.QueryEscape(id), Path: "/", HttpOnly: true, MaxAge: int(sm.MaxLife)}
	http.SetCookie(w, c)
	return s
}

type TestSession struct {
	ID      string
	LastAcs time.Time
	Val     map[string]interface{}
	Pvd     Provider
}

func (s *TestSession) Set(key string, val interface{}) error {
	s.Val[key] = val
	s.Pvd.SessionUpdate(s.ID)
	return nil
}

func (s *TestSession) Get(key string) interface{} {
	s.Pvd.SessionUpdate(s.ID)
	if v, ok := s.Val[key]; ok {
		return v
	}
	return nil
}

func (s *TestSession) Delete(key string) {
	delete(s.Val, key)
	s.Pvd.SessionUpdate(s.ID)
}

func (s *TestSession) SessionID() string {
	return s.ID
}

type TestProvider struct {
	Mu       sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func NewTestProvider() *TestProvider {
	ss := make(map[string]*list.Element)
	return &TestProvider{sessions: ss, list: list.New()}
}

func (p *TestProvider) SessionCreate(id string) (Session, error) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	v := make(map[string]interface{}, 0)
	s := &TestSession{ID: id, LastAcs: time.Now(), Val: v, Pvd: p}
	elem := p.list.PushBack(s)
	p.sessions[id] = elem
	return s, nil
}

func (p *TestProvider) SessionRead(id string) (Session, error) {
	if elem, ok := p.sessions[id]; ok {
		return elem.Value.(*TestSession), nil
	}
	return nil, fmt.Errorf(`session %q does not exist`, id)
}

func (p *TestProvider) SessionDestroy(id string) {
	if elem, ok := p.sessions[id]; ok {
		delete(p.sessions, id)
		p.list.Remove(elem)
	}
}

func (p *TestProvider) SessionGC(maxLife int64) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for elem := p.list.Back(); elem != nil; elem = p.list.Back() {
		if (elem.Value.(*TestSession).LastAcs.Unix() + maxLife) < time.Now().Unix() {
			p.list.Remove(elem)
			delete(p.sessions, elem.Value.(*TestSession).ID)
		} else {
			continue
		}
	}
}

func (p *TestProvider) SessionUpdate(id string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if elem, ok := p.sessions[id]; ok {
		elem.Value.(*TestSession).LastAcs = time.Now()
		p.list.MoveToFront(elem)
	}
	return nil
}
