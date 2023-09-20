package actions

import (
	"context"
	"fmt"

	"github.com/alayton/papers"
)

// Invalidates the current token chain
func Logout(ctx context.Context, p *papers.Papers, userID int64, accessToken *papers.AccessToken, refreshToken *papers.RefreshToken) error {
	if err := p.Config.Storage.Tokens.InvalidateTokenChain(ctx, userID, accessToken.Chain); err != nil {
		return fmt.Errorf("%w: %v", papers.ErrStorageError, err)
	}

	p.Config.Storage.TokenCache.SetChainValidity(accessToken.Chain, false)
	p.Config.Storage.TokenCache.SetTokenValidity(accessToken.Token, false)

	if refreshToken != nil {
		p.Config.Storage.TokenCache.SetTokenValidity(refreshToken.Token, false)

		if refreshToken.Chain != accessToken.Chain {
			// This shouldn't happen, but..
			p.Config.Storage.Tokens.InvalidateTokenChain(ctx, userID, refreshToken.Chain)
		}
	}

	return nil
}
