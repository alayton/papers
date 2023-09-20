package actions

import (
	"context"
	"fmt"
	"image"

	"github.com/pquerna/otp/totp"

	"github.com/alayton/papers"
)

type TOTPSetupResult struct {
	Secret string
	QR     image.Image
}

// Generates a TOTP secret and QR code to start the TOTP enrollment process
func TOTPSetup(ctx context.Context, p *papers.Papers, user papers.User) (*TOTPSetupResult, error) {
	if len(user.GetTOTPSecret()) > 0 {
		return nil, papers.ErrTOTPAlreadySetup
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      p.Config.ApplicationName,
		AccountName: user.GetEmail(),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", papers.ErrTOTPGenerateError, err)
	}

	img, err := key.Image(p.Config.TOTPQRSize, p.Config.TOTPQRSize)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", papers.ErrTOTPQRError, err)
	}

	return &TOTPSetupResult{
		Secret: key.Secret(),
		QR:     img,
	}, nil
}
