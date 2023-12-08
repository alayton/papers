package actions

import (
	"context"
	"fmt"

	"github.com/alayton/papers"
	"github.com/alayton/papers/validators"
)

type RecoverFinishFields struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// Starts the password recovery flow
func RecoverFinish(ctx context.Context, p *papers.Papers, fields RecoverFinishFields) error {
	if err := validators.Password(p, fields.Password); err != nil {
		return err
	}
	if len(fields.Token) == 0 {
		return papers.ErrTokenNotFound
	}

	user, err := p.Config.Storage.Users.GetUserByRecoveryToken(ctx, fields.Token)
	if err == papers.ErrUserNotFound {
		return papers.ErrTokenNotFound
	} else if err != nil {
		return err
	}

	hash, err := papers.HashPassword(p, fields.Password)
	if err != nil {
		return err
	}

	user.SetPassword(string(hash))
	user.SetRecoveryToken("")

	if err := p.Config.Storage.Users.UpdateUser(ctx, user); err != nil {
		p.Logger.Error("failed to update user during recovery", "error", err)
		return fmt.Errorf("%w: problem updating password", papers.ErrStorageError)
	}

	return nil
}
