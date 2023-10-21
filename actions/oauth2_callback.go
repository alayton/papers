package actions

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/alayton/papers"
)

type OAuth2CallbackFields struct {
	Provider string
	Nonce    string
	State    string
	Code     string
	Remember bool
}

type OAuth2CallbackResult struct {
	User         papers.User
	AccessToken  *papers.AccessToken
	RefreshToken *papers.RefreshToken
}

// Handles the callback from an OAuth2 provider
func OAuth2Callback(ctx context.Context, p *papers.Papers, fields OAuth2CallbackFields) (*OAuth2CallbackResult, error) {
	provider, ok := p.Config.OAuth2Providers[fields.Provider]
	if !ok {
		return nil, papers.ErrOAuth2BadProvider
	}
	if len(fields.State) == 0 || fields.State != fields.Nonce {
		return nil, papers.ErrOAuth2BadState
	}

	token, err := provider.Config.Exchange(ctx, fields.Code)
	if err != nil {
		return nil, papers.ErrOAuth2ExchangeFailed
	}

	identity, err := provider.GetIdentity(ctx, provider.Config, token)
	if err != nil {
		p.Logger.Error("failed to get identity from oauth2 provider", "error", err)
		return nil, papers.ErrOAuth2IdentityFailed
	}

	now := time.Now()
	user, err := p.Config.Storage.Users.GetUserByOAuth2Identity(ctx, fields.Provider, identity.GetID())
	if err != nil && err != papers.ErrUserNotFound {
		p.Logger.Error("failed to get user by oauth2 identity", "error", err)
		return nil, papers.ErrStorageError
	} else if err == papers.ErrUserNotFound {
		user, err = p.Config.Storage.Users.GetUserByEmail(ctx, identity.GetEmail())
		if err != nil && err != papers.ErrUserNotFound {
			p.Logger.Error("failed to get user by email", "error", err)
			return nil, papers.ErrStorageError
		} else if err == papers.ErrUserNotFound {
			user = p.Config.Storage.Users.NewUser()
			user.SetEmail(identity.GetEmail())
			user.SetUsername(fmt.Sprint("user", rand.Intn(1000000000)))
			user.SetCreatedAt(now)
			user.SetLastLogin(now)

			if err := p.Config.Storage.Users.CreateUser(ctx, user); err != nil {
				p.Logger.Error("failed to create user during oauth2 callback", "error", err)
				return nil, papers.ErrRegistrationFailed
			}
		}

		// User was either found by email or created, but lacks an oauth2 identity
		if err := p.Config.Storage.Users.CreateOAuth2Identity(ctx, user, fields.Provider, identity.GetID()); err != nil {
			p.Logger.Error("failed to create oauth2 identity during oauth2 callback", "error", err)
			return nil, papers.ErrRegistrationFailed
		}
	}

	// This won't be triggered if the user was just created
	if !user.GetLastLogin().Equal(now) {
		user.SetLastLogin(now)
		if err := p.Config.Storage.Users.UpdateUser(ctx, user); err != nil {
			p.Logger.Error("failed to update user's last login oauth2 callback", "error", err)
			return nil, papers.ErrStorageError
		}
	}

	accessToken, err := p.NewAccessToken(ctx, user.GetID(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
	}

	var refreshToken *papers.RefreshToken
	if fields.Remember {
		refreshToken, err = p.NewRefreshToken(ctx, accessToken)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
		}
	}

	return &OAuth2CallbackResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
