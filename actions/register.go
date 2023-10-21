package actions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alayton/papers"
	"github.com/alayton/papers/utils"
	"github.com/alayton/papers/validators"
)

type RegisterFields struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Creates a new user from the given fields. Errors are returned with errors.Join, use a type assertion with actions.MultiError to call Unwrap() []error
func Register(ctx context.Context, p *papers.Papers, fields RegisterFields) (papers.User, error) {
	emailErr := validators.Email(p, fields.Email)
	usernameErr := validators.Username(p, fields.Username)
	passwordErr := validators.Password(p, fields.Password)
	err := errors.Join(emailErr, usernameErr, passwordErr)
	if err != nil {
		return nil, err
	}

	hash, err := papers.HashPassword(p, fields.Password)
	if err != nil {
		return nil, errors.Join(err)
	}

	now := time.Now()

	user := p.Config.Storage.Users.NewUser()
	user.SetEmail(fields.Email)
	user.SetUsername(fields.Username)
	user.SetPassword(string(hash))
	user.SetCreatedAt(now)
	user.SetLastLogin(now)

	if p.Config.RequireConfirmation {
		token := utils.NewToken()
		user.SetConfirmToken(token)

		if p.Config.Mailer.Mailer != nil {
			p.Config.Mailer.Mailer.SendMessage(ctx, p, papers.Message{
				Type: papers.MessageConfirmation,
				To:   []papers.Email{{Address: fields.Email, Name: fields.Username}},
				Data: map[string]interface{}{
					"token": token,
				},
			})
		}
	}

	if err := p.Config.Storage.Users.CreateUser(ctx, user); err != nil {
		p.Logger.Error("failed to create user during registration", "error", err)
		return nil, errors.Join(fmt.Errorf("%w: problem creating a new user", papers.ErrRegistrationFailed))
	}

	return user, nil
}
