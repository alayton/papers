package actions

import (
	"context"

	"github.com/alayton/papers"
	"github.com/alayton/papers/utils"
	"github.com/alayton/papers/validators"
)

type RecoverStartFields struct {
	Email string `json:"email"`
}

// Starts the password recovery flow
func Recover(ctx context.Context, p *papers.Papers, fields RecoverStartFields) error {
	err := validators.Email(p, fields.Email)
	if err != nil {
		return err
	}

	user, err := p.Config.Storage.Users.GetUserByEmail(ctx, fields.Email)
	if err == papers.ErrUserNotFound {
		return nil
	} else if err != nil {
		return err
	}

	token := utils.NewToken()
	user.SetRecoveryToken(token)

	if err := p.Config.Storage.Users.UpdateUser(ctx, user); err != nil {
		return err
	}

	if p.Config.Mailer.Mailer != nil {
		if err := p.Config.Mailer.Mailer.SendMessage(ctx, p, papers.Message{
			Type: papers.MessageRecovery,
			To:   []papers.Email{{Address: user.GetEmail(), Name: user.GetUsername()}},
			Data: map[string]interface{}{
				"token": token,
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
