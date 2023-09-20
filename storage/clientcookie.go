package storage

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"

	"github.com/alayton/papers"
)

// Creates a new ClientCookie instance
// hashKey verifies the authenticity of the cookie using HMAC and should be at least 32 bytes (preferably 32 or 64 bytes)
// blockKey encrypts the cookie using AES-128, AES-192, or AES-256, and must be 16, 24, or 32 bytes respectively. A nil blockKey disables encryption
func NewClientCookie(hashKey []byte, blockKey []byte) *ClientCookie {
	return &ClientCookie{
		cookie:   securecookie.New(hashKey, blockKey),
		Secure:   true,
		HttpOnly: true,
	}
}

type ClientCookie struct {
	cookie   *securecookie.SecureCookie
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

func (c ClientCookie) Read(name string, r *http.Request, value interface{}) error {
	cookie, err := r.Cookie(name)
	if err == http.ErrNoCookie {
		return papers.ErrCookieNotFound
	} else if err != nil {
		// No other errors should be returned, but it feels icky to not include this case
		return papers.ErrCookieError
	}

	if err := c.cookie.Decode(name, cookie.Value, &value); err != nil {
		return fmt.Errorf("%w: %v", papers.ErrCookieDecodeError, err)
	}

	return nil
}

func (c ClientCookie) Write(name string, w http.ResponseWriter, maxAge time.Duration, value interface{}) error {
	encoded, err := c.cookie.Encode(name, value)
	if err != nil {
		return fmt.Errorf("%w: %v", papers.ErrCookieEncodeError, err)
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     c.Path,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		Domain:   c.Domain,
		MaxAge:   int(maxAge / time.Second),
	}
	http.SetCookie(w, cookie)

	return nil
}

func (c ClientCookie) Remove(name string, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     c.Path,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		Domain:   c.Domain,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}
