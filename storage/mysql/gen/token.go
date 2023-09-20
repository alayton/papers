package gen

import (
	"time"
)

func (t *AccessToken) GetUserID() int64        { return int64(t.UserID) }
func (t *AccessToken) GetToken() string        { return t.Token }
func (t *AccessToken) GetValid() bool          { return t.Valid != 0 }
func (t *AccessToken) GetChain() string        { return t.Chain }
func (t *AccessToken) GetCreatedAt() time.Time { return t.CreatedAt }

func (t *AccessToken) SetUserID(id int64)    { t.UserID = int32(id) }
func (t *AccessToken) SetToken(token string) { t.Token = token }
func (t *AccessToken) SetValid(valid bool) {
	if valid {
		t.Valid = 1
	} else {
		t.Valid = 0
	}
}
func (t *AccessToken) SetChain(chain string)     { t.Chain = chain }
func (t *AccessToken) SetCreatedAt(at time.Time) { t.CreatedAt = at }

func (t *RefreshToken) GetUserID() int64        { return int64(t.UserID) }
func (t *RefreshToken) GetToken() string        { return t.Token }
func (t *RefreshToken) GetValid() bool          { return t.Valid != 0 }
func (t *RefreshToken) GetChain() string        { return t.Chain }
func (t *RefreshToken) GetCreatedAt() time.Time { return t.CreatedAt }

func (t *RefreshToken) SetUserID(id int64)    { t.UserID = int32(id) }
func (t *RefreshToken) SetToken(token string) { t.Token = token }
func (t *RefreshToken) SetValid(valid bool) {
	if valid {
		t.Valid = 1
	} else {
		t.Valid = 0
	}
}
func (t *RefreshToken) SetChain(chain string)     { t.Chain = chain }
func (t *RefreshToken) SetCreatedAt(at time.Time) { t.CreatedAt = at }
