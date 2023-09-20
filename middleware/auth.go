package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/alayton/papers"
)

func Auth(p *papers.Papers) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var accessToken *papers.AccessToken
			if err := p.Config.Storage.Client.Read(p.Config.AccessCookieName, r, &accessToken); err != nil && err != papers.ErrCookieNotFound {
				p.Logger.Print("Middleware access cookie error:", err)
			}

			if accessToken != nil {
				if accessToken.Expiration.Before(time.Now()) {
					// Token exists, but has expired
					accessToken = nil
				} else {
					// Token exists and is unexpired, now check if it's been invalidated
					if valid, found := p.Config.Storage.TokenCache.IsTokenValid(accessToken.Token, accessToken.Chain); !found {
						// Token wasn't found in cache, check the database
						if token, err := p.Config.Storage.Tokens.GetAccessToken(ctx, accessToken.Identity, accessToken.Token); err != nil {
							// Two possibilities: the token wasn't in the database (ErrTokenNotFound), or the query errored
							if err != papers.ErrTokenNotFound {
								// Query error, log it, assume it's invalid
								p.Logger.Print("Middleware access token storage error:", err)
								accessToken = nil
							} else {
								// Token wasn't in the database, whether it's valid or not depends on whether valid access tokens are supposed to be stored
								p.Config.Storage.TokenCache.SetTokenValidity(accessToken.Token, !p.Config.StoreAllAccessTokens)
							}
						} else if !token.GetValid() {
							// Token was in the database and was invalidated
							p.Config.Storage.TokenCache.SetTokenValidity(accessToken.Token, false)
							accessToken = nil
						}
					} else if !valid {
						// Token was in cache and is invalid
						accessToken = nil
					}
				}
			}

			if accessToken == nil || accessToken.Expiration.Before(time.Now()) {
				// Access token was not found, expired, or was invalid. Check for a refresh token to generate a new access token
				var refreshToken *papers.RefreshToken
				if err := p.Config.Storage.Client.Read(p.Config.RefreshCookieName, r, &refreshToken); err != nil && err != papers.ErrCookieNotFound {
					p.Logger.Print("Middleware refresh cookie error:", err)
				}

				if refreshToken != nil {
					// Refresh token exists, check if it's been invalidated
					if valid, found := p.Config.Storage.TokenCache.IsTokenValid(refreshToken.Token, refreshToken.Chain); !found {
						// Token wasn't found in cache, check the database
						if token, err := p.Config.Storage.Tokens.GetRefreshToken(ctx, refreshToken.Identity, refreshToken.Token); err != nil {
							// Two possibilities: the token wasn't in the database (ErrTokenNotFound), or the query errored. Either way, we treat it as invalid
							if err != papers.ErrTokenNotFound {
								// Query error, log it
								p.Logger.Print("Middleware refresh token storage error:", err)
							}
							p.Config.Storage.TokenCache.SetTokenValidity(refreshToken.Token, false)
							refreshToken = nil
						} else if !token.GetValid() {
							// Token was in the database and was invalidated
							p.Config.Storage.TokenCache.SetTokenValidity(refreshToken.Token, false)
							refreshToken = nil
						}
					} else if !valid {
						// Token was in cache and is invalid
						refreshToken = nil
					}
				}

				if refreshToken != nil {
					// Token wasn't invalidated, we can use it
					if access, refresh, err := p.NewAccessTokenFromRefreshToken(ctx, refreshToken); err != nil {
						p.Logger.Print("Middleware access token refresh error:", err)
					} else {
						if err := p.Config.Storage.Client.Write(p.Config.AccessCookieName, w, p.Config.AccessExpiration, access); err != nil {
							p.Logger.Print("Middleware cookie write error:", err)
						}
						if err := p.Config.Storage.Client.Write(p.Config.RefreshCookieName, w, p.Config.RefreshExpiration, refresh); err != nil {
							p.Logger.Print("Middleware cookie write error:", err)
						}
						accessToken = access
					}
				}
			}

			if accessToken != nil {
				// If we still have an access token, it's either valid or been refreshed, so grab the user and store it in context
				user, err := p.Config.Storage.Users.GetUserByID(ctx, accessToken.Identity)
				if err != papers.ErrUserNotFound {
					p.Logger.Print("Middleware user storage error:", err)
				} else if err == nil {
					r = r.WithContext(context.WithValue(ctx, p.Config.UserContextKey, user))
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
