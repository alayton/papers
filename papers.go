package papers

import (
	"context"
	"log/slog"
)

type Papers struct {
	Config Config
	Logger *slog.Logger
	Roles  map[string][]string
}

func New() *Papers {
	papers := &Papers{
		Logger: slog.Default(),
		Roles:  map[string][]string{},
	}
	papers.SetDefaultConfig()

	return papers
}

func (p *Papers) Start(ctx context.Context) error {
	if p.Config.Storage.Users == nil {
		return ErrNoUserStorage
	}
	if p.Config.Storage.Tokens == nil {
		return ErrNoTokenStorage
	}
	if p.Config.Storage.Cookies == nil {
		return ErrNoClientStorage
	}
	if p.Config.Storage.Session == nil {
		return ErrNoSessionStorage
	}
	if p.Config.RouteParams == nil {
		return ErrNoRouteParams
	}

	if p.Config.PruneAccessTokensInterval > 0 {
		go p.pruneAccessTokens(ctx)
	}
	if p.Config.PruneRefreshTokensInterval > 0 {
		go p.pruneRefreshTokens(ctx)
	}
	return nil
}
