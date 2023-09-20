package actions

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/pquerna/otp/totp"

	"github.com/alayton/papers"
)

type TOTPFinishFields struct {
	Secret string `json:"secret"`
	Code   string `json:"code"`
}

// Generates a TOTP secret and QR code to start the TOTP enrollment process
func TOTPFinish(ctx context.Context, p *papers.Papers, user papers.User, fields TOTPFinishFields) error {
	if len(user.GetTOTPSecret()) > 0 {
		return papers.ErrTOTPAlreadySetup
	} else if !totp.Validate(fields.Code, fields.Secret) {
		return papers.ErrTOTPMismatch
	}

	secret := fields.Secret
	if len(p.Config.TOTPSecretEncryptionKey) > 0 {
		block, err := aes.NewCipher([]byte(p.Config.TOTPSecretEncryptionKey))
		if err != nil {
			return fmt.Errorf("%w: %v", papers.ErrCryptoError, err)
		}

		blockSize := block.BlockSize()
		buf := make([]byte, blockSize+len(secret))
		iv := buf[:blockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return fmt.Errorf("%w: %v", papers.ErrCryptoError, err)
		}

		cfb := cipher.NewCFBEncrypter(block, iv)
		cfb.XORKeyStream(buf[blockSize:], []byte(secret))
		secret = base64.StdEncoding.EncodeToString(buf)
	}

	user.SetTOTPSecret(secret)
	if err := p.Config.Storage.Users.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("%w: %v", papers.ErrStorageError, err)
	}

	return nil
}
