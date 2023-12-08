package papers

type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

const (
	// Startup errors
	ErrNoUserStorage    = ConstError("No UserStorage defined")
	ErrNoTokenStorage   = ConstError("No TokenStorage defined")
	ErrNoClientStorage  = ConstError("No CookieStorage defined")
	ErrNoSessionStorage = ConstError("No SessionStorage defined")
	ErrNoRouteParams    = ConstError("No RouteParams defined")

	// Generic errors
	ErrStorageError = ConstError("Unexpected storage error")
	ErrCryptoError  = ConstError("Unexpected cryptography error")

	// Storage errors
	ErrUserNotFound      = ConstError("User not found")
	ErrTokenNotFound     = ConstError("Token not found")
	ErrCookieNotFound    = ConstError("Cookie not found")
	ErrCookieError       = ConstError("Unexpected cookie error")
	ErrCookieDecodeError = ConstError("Couldn't decode cookie")
	ErrCookieEncodeError = ConstError("Couldn't encode cookie")
	ErrSessionError      = ConstError("Unexpected session error")
	ErrSessionMissingKey = ConstError("Session doesn't contain this key")

	// Registration errors
	ErrRegistrationFailed = ConstError("Registration failed")
	ErrDuplicateEmail     = ConstError("Email already in use")
	ErrInvalidEmail       = ConstError("Invalid email address")
	ErrMissingEmail       = ConstError("Email is required")
	ErrDuplicateUsername  = ConstError("Username already in use")
	ErrInvalidUsername    = ConstError("Invalid username")
	ErrMissingUsername    = ConstError("Username is required")
	ErrUsernameTooShort   = ConstError("Username is too short")
	ErrInvalidPassword    = ConstError("Invalid password")
	ErrPasswordError      = ConstError("There was a problem with the password")

	// Login errors
	ErrPasswordMismatch = ConstError("Password mismatch")
	ErrLoginFailed      = ConstError("Login failed")
	ErrUserLocked       = ConstError("Account is locked")

	// TOTP errors
	ErrTOTPUnexpected    = ConstError("Unexpected TOTP request")
	ErrTOTPGenerateError = ConstError("Unexpected TOTP generation error")
	ErrTOTPAlreadySetup  = ConstError("User already has TOTP setup")
	ErrTOTPQRError       = ConstError("Failed to create TOTP QR code")
	ErrTOTPMismatch      = ConstError("TOTP code doesn't match")

	// OAuth2 errors
	ErrOAuth2BadProvider    = ConstError("Invalid OAuth2 provider")
	ErrOAuth2BadState       = ConstError("Invalid OAuth2 state")
	ErrOAuth2LoginFailed    = ConstError("OAuth2 login failed")
	ErrOAuth2ExchangeFailed = ConstError("OAuth2 token exchange failed")
	ErrOAuth2IdentityFailed = ConstError("Couldn't get OAuth2 identity")

	// Mailer errors
	ErrMessageFailed     = ConstError("Failed to send email")
	ErrNoMessageTemplate = ConstError("Missing email template")
)
