package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/alayton/papers"
)

type LoginFields struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type LoginResult struct {
	NeedsTOTP    bool
	Remember     bool
	User         papers.User
	AccessToken  *papers.AccessToken
	RefreshToken *papers.RefreshToken
}

// Attempts to authenticate a user with the given email and password. Caller is responsible for sending access and refresh tokens to the client
func Login(ctx context.Context, p *papers.Papers, fields LoginFields) (*LoginResult, error) {
	user, err := p.Config.Storage.Users.GetUserByEmail(ctx, fields.Email)
	if err == papers.ErrUserNotFound {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
	}

	now := time.Now()
	if now.Before(user.GetLockedUntil()) {
		return nil, papers.ErrUserLocked
	}

	hash := user.GetPassword()
	if err := papers.ComparePassword(hash, fields.Password); err != nil {
		locked, err := failedLoginAttempt(ctx, p, user)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
		} else if locked {
			return nil, papers.ErrUserLocked
		}
		return nil, papers.ErrPasswordMismatch
	}

	if p.UserHasTOTP(user) {
		return &LoginResult{NeedsTOTP: true, User: user, Remember: fields.Remember}, nil
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
	if fields.Remember {
		refreshToken, err = p.NewRefreshToken(ctx, accessToken)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", papers.ErrLoginFailed, err)
		}
	}

	return &LoginResult{User: user, AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Handles logic around a failed login attempt (either password or TOTP mismatch) and returns if the user is locked or an error
func failedLoginAttempt(ctx context.Context, p *papers.Papers, user papers.User) (bool, error) {
	now := time.Now()
	lastAttempt := user.GetLastAttempt()
	if lastAttempt.IsZero() || lastAttempt.Add(p.Config.LockWindow).Before(now) {
		user.SetLastAttempt(now)
		user.SetAttempts(1)
	} else {
		user.SetAttempts(user.GetAttempts() + 1)
	}

	if user.GetAttempts() >= p.Config.LockAttempts {
		user.SetLockedUntil(now.Add(p.Config.LockDuration))
	}

	if err := p.Config.Storage.Users.UpdateUser(ctx, user); err != nil {
		return true, err
	}

	// GetLockedUntil will either have a valid date or a zero date if the account has never been locked, both of which are valid for this comparison
	return now.Before(user.GetLockedUntil()), nil
}
