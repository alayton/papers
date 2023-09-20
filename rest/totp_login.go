package rest

import (
	"errors"
	"net/http"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

type totpLoginResponse struct {
	User  papers.User `json:"user,omitempty"`
	Error string      `json:"error,omitempty"`
}

func TOTPLogin(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp totpLoginResponse

		var fields actions.TOTPLoginFields
		if err := readJSON(r.Body, &fields); err != nil {
			resp.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, resp)
			return
		}

		result, err := actions.TOTPLogin(r.Context(), p, fields)
		if err != nil {
			status := http.StatusBadRequest
			if err == papers.ErrUserNotFound {
				resp.Error = "User not found"
			} else if err == papers.ErrTOTPMismatch {
				resp.Error = "Incorrect TOTP code"
			} else if errors.Is(err, papers.ErrLoginFailed) {
				resp.Error = papers.ErrLoginFailed.Error()
				status = http.StatusInternalServerError
			} else {
				resp.Error = err.Error()
			}

			w.WriteHeader(status)
			writeJSON(w, resp)
			return
		}

		status := http.StatusOK

		if result.RefreshToken != nil {
			if err := p.Config.Storage.Client.Write(p.Config.RefreshCookieName, w, p.Config.RefreshExpiration, result.RefreshToken); err != nil {
				status = http.StatusInternalServerError
				resp.Error = err.Error()
			}
		}

		if err := p.Config.Storage.Client.Write(p.Config.AccessCookieName, w, p.Config.AccessExpiration, result.AccessToken); err != nil {
			status = http.StatusInternalServerError
			resp.Error = err.Error()
		} else if status == http.StatusOK {
			resp.User = result.User
		}

		w.WriteHeader(status)
		writeJSON(w, resp)
	}
}
