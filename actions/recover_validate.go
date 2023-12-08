package actions

import (
	"context"

	"github.com/alayton/papers"
)

type RecoverValidateFields struct {
	Token string
}

// Starts the password recovery flow
func RecoverValidate(ctx context.Context, p *papers.Papers, fields RecoverValidateFields) error {
	if len(fields.Token) == 0 {
		return papers.ErrTokenNotFound
	}

	_, err := p.Config.Storage.Users.GetUserByRecoveryToken(ctx, fields.Token)
	if err == papers.ErrUserNotFound {
		return papers.ErrTokenNotFound
	} else if err != nil {
		return err
	}

	return nil
}
