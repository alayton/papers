package rest

import (
	"errors"
	"net/http"
	"strconv"

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

		strID, err := p.Config.Storage.Session.Get(r, "uid")
		if err != nil {
			if err == papers.ErrSessionMissingKey {
				resp.Error = err.Error()
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				resp.Error = papers.ErrSessionError.Error()
				w.WriteHeader(http.StatusInternalServerError)
			}
			writeJSON(w, resp)
			return
		}

		userID, err := strconv.ParseInt(strID, 10, 64)
		if err != nil {
			resp.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, resp)
			return
		}

		remember, _ := p.Config.Storage.Session.Get(r, "rm")

		var fields actions.TOTPLoginFields
		if err := readJSON(r.Body, &fields); err != nil {
			resp.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, resp)
			return
		}

		fields.UserID = userID
		fields.Remember = remember == "true"

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
			if err := p.Config.Storage.Cookies.Write(p.Config.RefreshCookieName, w, p.Config.RefreshExpiration, result.RefreshToken); err != nil {
				status = http.StatusInternalServerError
				resp.Error = err.Error()
			}
		}

		if err := p.Config.Storage.Cookies.Write(p.Config.AccessCookieName, w, p.Config.AccessExpiration, result.AccessToken); err != nil {
			status = http.StatusInternalServerError
			resp.Error = err.Error()
		} else if status == http.StatusOK {
			resp.User = result.User

			p.Config.Storage.Session.MultiDelete(r, []string{"uid", "rm"})
			p.Config.Storage.Session.Write(r, w)
		}

		w.WriteHeader(status)
		writeJSON(w, resp)
	}
}
