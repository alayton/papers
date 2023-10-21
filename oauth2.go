package papers

import (
	"context"

	"golang.org/x/oauth2"
)

type OAuth2Provider struct {
	Config      oauth2.Config
	GetIdentity func(ctx context.Context, cfg oauth2.Config, token *oauth2.Token) (OAuth2Identity, error)
}

type OAuth2Identity interface {
	GetID() string
	GetEmail() string
}
