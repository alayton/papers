package papers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p *Papers, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), p.Config.BCryptCost)
	if err != nil {
		p.Logger.Print("Error from bcrypt.GenerateFromPassword:", err)
		return "", ErrPasswordError
	}
	return string(hash), nil
}

func ComparePassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrInvalidPassword
	}
	return nil
}
