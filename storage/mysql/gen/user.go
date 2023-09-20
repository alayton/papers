package gen

import (
	"database/sql"
	"time"
)

func (u *User) GetID() int64              { return u.ID }
func (u *User) GetEmail() string          { return u.Email }
func (u *User) GetUsername() string       { return u.Username.String }
func (u *User) GetPassword() string       { return u.Password.String }
func (u *User) GetConfirmed() bool        { return u.Confirmed != 0 }
func (u *User) GetConfirmToken() string   { return u.ConfirmToken.String }
func (u *User) GetRecoveryToken() string  { return u.RecoveryToken.String }
func (u *User) GetLockedUntil() time.Time { return u.LockedUntil.Time }
func (u *User) GetAttempts() int          { return int(u.Attempts) }
func (u *User) GetLastAttempt() time.Time { return u.LastAttempt.Time }
func (u *User) GetCreatedAt() time.Time   { return u.CreatedAt }
func (u *User) GetLastLogin() time.Time   { return u.LastLogin }
func (u *User) GetTOTPSecret() string     { return u.TotpSecret.String }

func (u *User) SetID(id int64)        { u.ID = id }
func (u *User) SetEmail(email string) { u.Email = email }
func (u *User) SetUsername(username string) {
	if len(username) == 0 {
		u.Username = sql.NullString{}
	} else {
		u.Username = sql.NullString{String: username, Valid: true}
	}
}
func (u *User) SetPassword(password string) {
	if len(password) == 0 {
		u.Password = sql.NullString{}
	} else {
		u.Password = sql.NullString{String: password, Valid: true}
	}
}
func (u *User) SetConfirmed(confirmed bool) {
	if confirmed {
		u.Confirmed = 1
	} else {
		u.Confirmed = 0
	}
}
func (u *User) SetConfirmToken(token string) {
	if len(token) == 0 {
		u.ConfirmToken = sql.NullString{}
	} else {
		u.ConfirmToken = sql.NullString{String: token, Valid: true}
	}
}
func (u *User) SetRecoveryToken(token string) {
	if len(token) == 0 {
		u.RecoveryToken = sql.NullString{}
	} else {
		u.RecoveryToken = sql.NullString{String: token, Valid: true}
	}
}
func (u *User) SetLockedUntil(until time.Time) {
	if until.IsZero() {
		u.LockedUntil = sql.NullTime{}
	} else {
		u.LockedUntil = sql.NullTime{Time: until, Valid: true}
	}
}
func (u *User) SetAttempts(attempts int) { u.Attempts = int32(attempts) }
func (u *User) SetLastAttempt(at time.Time) {
	if at.IsZero() {
		u.LastAttempt = sql.NullTime{}
	} else {
		u.LastAttempt = sql.NullTime{Time: at, Valid: true}
	}
}
func (u *User) SetCreatedAt(at time.Time) { u.CreatedAt = at }
func (u *User) SetLastLogin(at time.Time) { u.LastLogin = at }
func (u *User) SetTOTPSecret(secret string) {
	if len(secret) == 0 {
		u.TotpSecret = sql.NullString{}
	} else {
		u.TotpSecret = sql.NullString{String: secret, Valid: true}
	}
}
