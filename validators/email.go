package validators

import (
	"net/mail"

	"github.com/alayton/papers"
)

func Email(p *papers.Papers, address string) error {
	if len(address) == 0 {
		return papers.ErrMissingEmail
	}

	email, err := mail.ParseAddress(address)
	if err != nil || email.Address != address {
		return papers.ErrInvalidEmail
	}
	return nil
}
