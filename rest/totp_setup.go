package rest

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"net/http"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

type totpSetupResponse struct {
	Secret string `json:"secret,omitempty"`
	QR     string `json:"qr,omitempty"`
	Error  string `json:"error,omitempty"`
}

func TOTPSetup(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, loggedIn := p.LoggedInUser(r)
		if !loggedIn {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var resp totpSetupResponse

		result, err := actions.TOTPSetup(r.Context(), p, user)
		if err != nil {
			status := http.StatusInternalServerError
			if err == papers.ErrTOTPAlreadySetup {
				status = http.StatusBadRequest
			}

			e := errors.Unwrap(err)
			if e != nil {
				resp.Error = e.Error()
			} else {
				resp.Error = err.Error()
			}

			w.WriteHeader(status)
			writeJSON(w, resp)
			return
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, result.QR); err != nil {
			p.Logger.Print("Failed to encode TOTP QR code:", err)
			resp.Error = papers.ErrTOTPQRError.Error()
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, resp)
			return
		}

		resp.Secret = result.Secret
		resp.QR = base64.StdEncoding.EncodeToString(buf.Bytes())

		writeJSON(w, resp)
	}
}
