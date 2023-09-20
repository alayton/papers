package papers

import (
	"net/http"
	"time"
)

type User interface {
	GetID() int64
	GetEmail() string
	GetUsername() string
	GetPassword() string
	GetConfirmed() bool
	GetConfirmToken() string
	GetRecoveryToken() string
	GetLockedUntil() time.Time
	GetAttempts() int
	GetLastAttempt() time.Time
	GetCreatedAt() time.Time
	GetLastLogin() time.Time
	GetTOTPSecret() string

	SetID(id int64)
	SetEmail(email string)
	SetUsername(username string)
	SetPassword(password string)
	SetConfirmed(confirmed bool)
	SetConfirmToken(token string)
	SetRecoveryToken(token string)
	SetLockedUntil(until time.Time)
	SetAttempts(attempts int)
	SetLastAttempt(at time.Time)
	SetCreatedAt(at time.Time)
	SetLastLogin(at time.Time)
	SetTOTPSecret(secret string)
}

func (p *Papers) LoggedInUser(r *http.Request) (User, bool) {
	user, ok := r.Context().Value(p.Config.UserContextKey).(User)
	return user, ok
}

func (p *Papers) UserHasTOTP(u User) bool {
	return len(u.GetTOTPSecret()) > 0
}

func (p *Papers) IsUserConfirmed(u User) bool {
	if p.Config.RequireConfirmation && !u.GetConfirmed() {
		return false
	}
	return true
}

func (p *Papers) IsUserLocked(u User) bool {
	if p.Config.Locking {
		now := time.Now()
		unlockTime := now.Add(p.Config.LockDuration)
		if now.Before(unlockTime) {
			return true
		}
	}
	return false
}

func (p *Papers) IsUserValid(u User) bool {
	if p.Config.RequireConfirmation && !u.GetConfirmed() {
		return false
	} else if p.Config.Locking {
		now := time.Now()
		unlockTime := now.Add(p.Config.LockDuration)
		if now.Before(unlockTime) {
			return false
		}
	}
	return true
}
