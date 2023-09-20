package papers

import (
	"context"
	"net/http"
	"time"
)

type UserStorage interface {
	// Factory method that returns an empty User
	NewUser() User
	// Persists a new user. The storage implementation must call User.SetID with the new ID
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	CreateOAuth2Identity(ctx context.Context, user User, provider, identity string) error
	RemoveOAuth2Identity(ctx context.Context, user User, provider, identity string) error
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserByOAuth2Identity(ctx context.Context, provider, identity string) (User, error)
	GetUserByConfirmationToken(ctx context.Context, token string) (User, error)
	GetUserByRecoveryToken(ctx context.Context, token string) (User, error)
	GetUserRoles(ctx context.Context, user User) ([]string, error)
	GetUserPermissions(ctx context.Context, user User) ([]string, error)
}

type TokenCache interface {
	IsTokenValid(token, chain string) (bool, bool)
	SetTokenValidity(token string, valid bool)
	SetChainValidity(chain string, valid bool)
}

type TokenStorage interface {
	CreateAccessToken(ctx context.Context, userID int64, token, chain string, valid bool) error
	CreateRefreshToken(ctx context.Context, userID int64, token, chain string, valid bool) error
	GetAccessToken(ctx context.Context, userID int64, token string) (Token, error)
	GetRefreshToken(ctx context.Context, userID int64, token string) (Token, error)
	InvalidateAccessTokens(ctx context.Context, userID int64) error
	InvalidateRefreshToken(ctx context.Context, userID int64, token string) error
	InvalidateRefreshTokens(ctx context.Context, userID int64) error
	InvalidateTokenChain(ctx context.Context, userID int64, chain string) error
	PruneAccessTokens(ctx context.Context, timeToStale time.Duration) error
	PruneRefreshTokens(ctx context.Context, timeToStale time.Duration) error
}

type ClientStorage interface {
	Read(name string, r *http.Request, value interface{}) error
	Write(name string, w http.ResponseWriter, maxAge time.Duration, value interface{}) error
	Remove(name string, w http.ResponseWriter)
}

type TTLStorage interface {
	Set(key string, value string, expiration time.Duration)
	Get(key string) (string, bool)
}
