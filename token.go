package papers

import (
	"context"
	"fmt"
	"time"

	"github.com/alayton/papers/utils"
)

type Token interface {
	GetUserID() int64
	GetToken() string
	GetValid() bool
	GetChain() string
	GetCreatedAt() time.Time

	SetUserID(id int64)
	SetToken(token string)
	SetValid(valid bool)
	SetChain(chain string)
	SetCreatedAt(at time.Time)
}

type AccessToken struct {
	Identity    int64     `json:"id"`
	Token       string    `json:"to"`
	Expiration  time.Time `json:"ex"`
	Roles       []string  `json:"ro"`
	Permissions []string  `json:"pe"`
	Chain       string    `json:"ch"`
}

type RefreshToken struct {
	Identity int64  `json:"id"`
	Token    string `json:"to"`
	Chain    string `json:"ch"`
	Limited  bool   `json:"li"`
}

func (p *Papers) NewAccessToken(ctx context.Context, userID int64, refresh *RefreshToken) (*AccessToken, error) {
	var chain string
	if refresh != nil {
		chain = refresh.Chain
	} else {
		chain = utils.NewChain()
	}

	token := &AccessToken{
		Identity:   userID,
		Token:      utils.NewToken(),
		Expiration: time.Now().Add(p.Config.AccessExpiration),
		Chain:      chain,
	}

	if p.Config.StoreAllAccessTokens {
		if err := p.Config.Storage.Tokens.CreateAccessToken(ctx, userID, token.Token, token.Chain, true); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrStorageError, err)
		}
	}

	p.Config.Storage.TokenCache.SetTokenValidity(token.Token, true)

	return token, nil
}

func (p *Papers) NewRefreshToken(ctx context.Context, access *AccessToken) (*RefreshToken, error) {
	token := &RefreshToken{
		Identity: access.Identity,
		Token:    utils.NewToken(),
		Chain:    access.Chain,
	}

	if err := p.Config.Storage.Tokens.CreateRefreshToken(ctx, access.Identity, token.Token, token.Chain, true); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStorageError, err)
	}

	p.Config.Storage.TokenCache.SetTokenValidity(token.Token, true)

	return token, nil
}

func (p *Papers) NewAccessTokenFromRefreshToken(ctx context.Context, refresh *RefreshToken) (*AccessToken, *RefreshToken, error) {
	access, err := p.NewAccessToken(ctx, refresh.Identity, refresh)
	if err != nil {
		return nil, nil, err
	}

	if p.Config.RotateRefreshTokens {
		if err := p.Config.Storage.Tokens.InvalidateRefreshToken(ctx, refresh.Identity, refresh.Token); err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrStorageError, err)
		}
		p.Config.Storage.TokenCache.SetTokenValidity(refresh.Token, false)

		refresh, err = p.NewRefreshToken(ctx, access)
		if err != nil {
			return nil, nil, err
		}
	}

	return access, refresh, nil
}

func (p *Papers) pruneAccessTokens(ctx context.Context) {
	if err := p.Config.Storage.Tokens.PruneAccessTokens(ctx, p.Config.StaleAccessTokensAge); err != nil {
		p.Logger.Print("PruneAccessTokens error:", err)
	}

	timer := time.NewTimer(p.Config.PruneAccessTokensInterval)
	select {
	case <-ctx.Done():
		timer.Stop()
		return
	case <-timer.C:
		p.pruneAccessTokens(ctx)
	}
}

func (p *Papers) pruneRefreshTokens(ctx context.Context) {
	if err := p.Config.Storage.Tokens.PruneRefreshTokens(ctx, p.Config.StaleRefreshTokensAge); err != nil {
		p.Logger.Print("PruneRefreshTokens error:", err)
	}

	timer := time.NewTimer(p.Config.PruneRefreshTokensInterval)
	select {
	case <-ctx.Done():
		timer.Stop()
		return
	case <-timer.C:
		p.pruneRefreshTokens(ctx)
	}
}
