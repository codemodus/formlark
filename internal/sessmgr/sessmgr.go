package sessmgr

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Manager struct {
	name    string
	maxLife int64
	prov    Provider
}

func New(cookieName string, maxLife int64, provider Provider) *Manager {
	return &Manager{
		name:    cookieName,
		maxLife: maxLife,
		prov:    provider,
	}
}

func (m *Manager) genSessID() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", errors.New("bad rng")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (m *Manager) SessStart(w http.ResponseWriter, r *http.Request) (s *Session, err error) {
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

	id, err := m.genSessID()
	if err != nil {
		return nil, err
	}

	if s, err = m.prov.Create(id); err != nil {
		return nil, err
	}

	c = &http.Cookie{
		Name:     m.name,
		Value:    url.QueryEscape(id),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(m.maxLife),
	}
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
	if err != nil || id == "" {
		return
	}
	m.prov.Destroy(id)
}

func unsetCookie(w http.ResponseWriter, name string) {
	c := &http.Cookie{
		Name:    name,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, c)
}
