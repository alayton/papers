package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/alayton/papers"
)

var usernameCharacters = regexp.MustCompile(`^[a-zA-Z0-9_\.,<>/\?:;'"|\[\]\{\}\+=\(\)\*&\^%\$#@!~-]+$`)

func Username(p *papers.Papers, username string) error {
	length := len(username)
	if length == 0 && !p.Config.RequireUsername {
		return nil
	} else if length == 0 {
		return papers.ErrMissingUsername
	} else if length < p.Config.UsernameMinLength {
		return fmt.Errorf("%w: must be at least %d characters", papers.ErrInvalidUsername, p.Config.UsernameMinLength)
	} else if !usernameCharacters.MatchString(username) {
		return fmt.Errorf("%w: contains invalid characters", papers.ErrInvalidUsername)
	} else if p.Config.UniqueUsernames {
		_, err := p.Config.Storage.Users.GetUserByUsername(context.Background(), username)
		if err != papers.ErrUserNotFound {
			return papers.ErrDuplicateUsername
		}
	}
	return nil
}
