package actions

import (
	"context"
	"fmt"

	"github.com/alayton/papers"
)

type ConfirmFields struct {
	Token string
}

// Starts the password recovery flow
func Confirm(ctx context.Context, p *papers.Papers, fields ConfirmFields) error {
	if len(fields.Token) == 0 {
		return papers.ErrTokenNotFound
	}

	user, err := p.Config.Storage.Users.GetUserByConfirmationToken(ctx, fields.Token)
	if err == papers.ErrUserNotFound {
		return papers.ErrTokenNotFound
	} else if err != nil {
		return err
	}

	user.SetConfirmed(true)
	user.SetConfirmToken("")

	if err := p.Config.Storage.Users.UpdateUser(ctx, user); err != nil {
		p.Logger.Error("failed to update user during confirmation", "error", err)
		return fmt.Errorf("%w: problem confirming user", papers.ErrStorageError)
	}

	return nil
}
