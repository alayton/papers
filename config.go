package papers

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Storage struct {
		Users      UserStorage
		Tokens     TokenStorage
		Cookies    CookieStorage
		Session    SessionStorage
		TokenCache TokenCache
	}

	Mailer struct {
		Mailer Mailer
		From   Email
	}

	Routes struct {
		// Path of the confirmation page with a printf style placeholder for the token, e.g. /confirm/%s
		Confirm string
		// Path of the recovery page for the forgot password flow with a printf style placeholder for the token, e.g. /recover/%s
		Recovery string
	}

	// Complete root level URL of the application, with no trailing slash, e.g. https://example.com
	BaseURL string
	// Name of the application (used when generating TOTP secret)
	ApplicationName string

	// Key to use when storing the logged in User in the request context
	UserContextKey string

	// Require a standard email/password account to be confirmed before being usable
	RequireConfirmation bool
	// Require a username when registering an account
	RequireUsername bool
	// Require usernames to be unique
	UniqueUsernames bool
	// Minimum length of a username
	UsernameMinLength int

	BCryptCost        int
	PasswordMinLength int
	// Require both lower and upper case letters
	PasswordRequireMixedCase bool
	// Require at least one number
	PasswordRequireNumbers bool
	// Require at least one special character
	PasswordRequireSpecials bool
	// Password length that other character requirements are ignored. Intended to allow for passphrases
	PasswordRelaxedLength int

	// Width/height of TOTP setup QR code
	TOTPQRSize int
	// If set, key used to encrypt the TOTP secret before saving it in storage. Must be 16/24/32 bytes for AES-128/AES-192/AES-256 respectively
	TOTPSecretEncryptionKey string

	// Name of the login token cookie
	LoginCookieName string
	// Name of the access token cookie
	AccessCookieName string
	// Name of the refresh token cookie
	RefreshCookieName string

	// How long until an access token expires
	AccessExpiration time.Duration
	// How long until a refresh token expires
	RefreshExpiration time.Duration
	// How long until a recovery token expires
	RecoveryExpiration time.Duration
	// How long tokens are cached before expiring
	TokenCacheExpiration time.Duration

	// How much leeway is given when checking token expiration
	ExpirationLeeway time.Duration

	// Store all access tokens for extra verification, instead of only storing invalidated access tokens
	StoreAllAccessTokens bool
	// Invalidate and issue a new refresh token each time a refresh token is used
	RotateRefreshTokens bool

	// How long to wait between pruning expired access tokens. 0 to disable
	PruneAccessTokensInterval time.Duration
	// How long to wait between pruning expired refresh tokens. 0 to disable
	PruneRefreshTokensInterval time.Duration

	// How old access tokens must be before they are pruned
	StaleAccessTokensAge time.Duration
	// How old refresh tokens must be before they are pruned
	StaleRefreshTokensAge time.Duration

	// Temporarily lock accounts that have too many failed login attempts
	Locking bool
	// The number of failed login attempts to allow before locking an account
	LockAttempts int
	// How long to wait after a failed login attempt before resetting the number of attempts
	LockWindow time.Duration
	// How long to lock an account after too many failed attempts
	LockDuration time.Duration
}

func (p *Papers) SetDefaultConfig() {
	p.Config = Config{
		UserContextKey: "papers_user",

		RequireConfirmation: true,
		RequireUsername:     true,
		UniqueUsernames:     true,
		UsernameMinLength:   3,

		BCryptCost:               bcrypt.DefaultCost,
		PasswordMinLength:        12,
		PasswordRequireMixedCase: true,
		PasswordRequireNumbers:   true,
		PasswordRequireSpecials:  false,
		PasswordRelaxedLength:    20,

		TOTPQRSize: 200,

		LoginCookieName:   "login",
		AccessCookieName:  "access",
		RefreshCookieName: "refresh",

		AccessExpiration:     time.Minute * 15,
		RefreshExpiration:    time.Hour * 24 * 30,
		RecoveryExpiration:   time.Minute * 30,
		TokenCacheExpiration: time.Hour,

		ExpirationLeeway: time.Minute,

		StoreAllAccessTokens: false,
		RotateRefreshTokens:  true,

		PruneAccessTokensInterval:  time.Hour * 3,
		PruneRefreshTokensInterval: time.Hour * 24,

		StaleAccessTokensAge:  time.Hour,
		StaleRefreshTokensAge: time.Hour * 48,

		Locking:      true,
		LockAttempts: 4,
		LockWindow:   time.Minute,
		LockDuration: time.Minute * 5,
	}
}
