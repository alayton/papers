package actions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/alayton/papers"
)

type OAuth2StartFields struct {
	Provider string `json:"provider"`
	Remember bool   `json:"remember"`
}

type OAuth2StartResult struct {
	Nonce       string
	RedirectURL string
}

// Starts the OAuth2 flow for a given provider
func OAuth2Start(ctx context.Context, p *papers.Papers, fields OAuth2StartFields) (*OAuth2StartResult, error) {
	provider, ok := p.Config.OAuth2Providers[fields.Provider]
	if !ok {
		return nil, papers.ErrOAuth2BadProvider
	}

	nonceBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		p.Logger.Error("failed to read random bytes during oauth2 start", "error", err)
		return nil, papers.ErrCryptoError
	}
	nonce := base64.URLEncoding.EncodeToString(nonceBytes)

	url := provider.Config.AuthCodeURL(nonce)

	return &OAuth2StartResult{Nonce: nonce, RedirectURL: url}, nil
}
