package gorilla

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"github.com/alayton/papers"
)

// Creates a new Session instance
// hashKey verifies the authenticity of the cookie using HMAC and should be at least 32 bytes (preferably 32 or 64 bytes)
// blockKey encrypts the cookie using AES-128, AES-192, or AES-256, and must be 16, 24, or 32 bytes respectively. A nil blockKey disables encryption
func NewSession(name string, contextKey string, maxAge time.Duration, keyPairs ...[]byte) *Session {
	store := sessions.NewCookieStore(keyPairs...)
	store.MaxAge(int(maxAge / time.Second))
	store.Options.HttpOnly = true
	store.Options.Secure = true

	return &Session{
		Name:       name,
		ContextKey: contextKey,
		store:      store,
	}
}

type Session struct {
	Name       string
	ContextKey string
	store      sessions.Store
}

func (s Session) Get(r *http.Request, key string) (string, error) {
	session, err := s.store.Get(r, s.Name)
	if err != nil {
		if e, ok := err.(securecookie.Error); ok && !e.IsDecode() {
			return "", err
		}
		session, _ = s.store.New(r, s.Name)
	}

	v, ok := session.Values[key]
	if !ok {
		// Not typically a problem, simply can't use the (string, bool) pattern because of the session error handling above
		return "", papers.ErrSessionMissingKey
	}
	return v.(string), nil
}

func (s Session) Set(r *http.Request, key, value string) error {
	session, err := s.store.Get(r, s.Name)
	if err != nil {
		if e, ok := err.(securecookie.Error); ok && !e.IsDecode() {
			return err
		}
		session, _ = s.store.New(r, s.Name)
	}

	session.Values[key] = value
	return nil
}

func (s Session) MultiSet(r *http.Request, vals map[string]string) error {
	session, err := s.store.Get(r, s.Name)
	if err != nil {
		if e, ok := err.(securecookie.Error); ok && !e.IsDecode() {
			return err
		}
		session, _ = s.store.New(r, s.Name)
	}

	for k, v := range vals {
		session.Values[k] = v
	}
	return nil
}

func (s Session) Delete(r *http.Request, key string) error {
	session, err := s.store.Get(r, s.Name)
	if err != nil {
		if e, ok := err.(securecookie.Error); ok && !e.IsDecode() {
			return err
		}
		session, _ = s.store.New(r, s.Name)
	}

	delete(session.Values, key)
	return nil
}

func (s Session) MultiDelete(r *http.Request, keys []string) error {
	session, err := s.store.Get(r, s.Name)
	if err != nil {
		if e, ok := err.(securecookie.Error); ok && !e.IsDecode() {
			return err
		}
		session, _ = s.store.New(r, s.Name)
	}

	for _, k := range keys {
		delete(session.Values, k)
	}
	return nil
}

func (s Session) Write(r *http.Request, w http.ResponseWriter) error {
	session, err := s.store.Get(r, s.Name)
	if err != nil {
		if e, ok := err.(securecookie.Error); ok && !e.IsDecode() {
			return err
		}
		session, _ = s.store.New(r, s.Name)
	}

	return s.store.Save(r, w, session)
}
