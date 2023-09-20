package validators

import (
	"fmt"
	"regexp"

	"github.com/alayton/papers"
)

var passwordLowerCase = regexp.MustCompile(`[a-z]`)
var passwordUpperCase = regexp.MustCompile(`[A-Z]`)
var passwordNumber = regexp.MustCompile(`[0-9]`)
var passwordSpecial = regexp.MustCompile(`[!@#$%^&*+=_-]`)

func Password(p *papers.Papers, password string) error {
	bytes := []byte(password)
	if len(bytes) > 72 {
		// bcrypt enforces a limit of 72 bytes (not characters, utf8 encoding), i choose to limit input length rather than truncate
		return fmt.Errorf("%w: must be no longer than 72 characters", papers.ErrInvalidPassword)
	}

	length := len(password)
	if length < p.Config.PasswordMinLength {
		return fmt.Errorf("%w: must be at least %d characters", papers.ErrInvalidPassword)
	} else if length >= p.Config.PasswordRelaxedLength {
		return nil
	}

	lowercase := passwordLowerCase.MatchString(password)
	uppercase := passwordUpperCase.MatchString(password)
	number := passwordNumber.MatchString(password)
	special := passwordSpecial.MatchString(password)
	if p.Config.PasswordRequireMixedCase && (!lowercase || !uppercase) {
		return fmt.Errorf("%w: must have both lower and upper case characters", papers.ErrInvalidPassword)
	} else if p.Config.PasswordRequireNumbers && !number {
		return fmt.Errorf("%w: must contain at least one number", papers.ErrInvalidPassword)
	} else if p.Config.PasswordRequireSpecials && !special {
		return fmt.Errorf("%w: must contain at least one special character (!@#$%^&*+=_-)", papers.ErrInvalidPassword)
	}
	return nil
}
