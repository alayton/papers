package actions

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/alayton/papers"
)

type TOTPLoginFields struct {
	LoginToken *papers.LoginToken
	Code       string `json:"code"`
}

type TOTPLoginResult struct {
	User         papers.User
	AccessToken  *papers.AccessToken
	RefreshToken *papers.RefreshToken
}

// Completes authentication of users with a TOTP configured. Caller is responsible for generating access and refresh tokens for the returned User
func TOTPLogin(ctx context.Context, p *papers.Papers, fields TOTPLoginFields) (*TOTPLoginResult, error) {
	now := time.Now()
	if now.After(fields.LoginToken.Expiration) {
		return nil, papers.ErrLoginTokenExpired
	}

	user, err := p.Config.Storage.Users.GetUserByID(ctx, fields.LoginToken.Identity)
	if err == papers.ErrUserNotFound {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
	}

	if now.Before(user.GetLockedUntil()) {
		return nil, papers.ErrUserLocked
	}

	secret := user.GetTOTPSecret()
	if len(p.Config.TOTPSecretEncryptionKey) > 0 {
		block, err := aes.NewCipher([]byte(p.Config.TOTPSecretEncryptionKey))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", papers.ErrCryptoError, err)
		}

		blockSize := block.BlockSize()
		decoded, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", papers.ErrCryptoError, err)
		} else if len(decoded) < blockSize {
			return nil, fmt.Errorf("%w: %s", papers.ErrCryptoError, "secret is too short")
		}

		iv := decoded[:blockSize]
		encrypted := decoded[blockSize:]

		cfb := cipher.NewCFBDecrypter(block, iv)
		cfb.XORKeyStream(encrypted, encrypted)
		secret = string(encrypted)
	}

	if !totp.Validate(fields.Code, secret) {
		locked, err := failedLoginAttempt(ctx, p, user)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
		} else if locked {
			return nil, papers.ErrUserLocked
		}
		return nil, papers.ErrTOTPMismatch
	}

	user.SetLastLogin(time.Now())
	if err := p.Config.Storage.Users.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
	}

	accessToken, err := p.NewAccessToken(ctx, user.GetID(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
	}

	var refreshToken *papers.RefreshToken
	if fields.LoginToken.Remember {
		refreshToken, err = p.NewRefreshToken(ctx, accessToken)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
		}
	}

	return &TOTPLoginResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
