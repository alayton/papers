package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

type loginResponse struct {
	NeedsTOTP bool        `json:"needsTotp,omitempty"`
	User      papers.User `json:"user,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func Login(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp loginResponse

		var fields actions.LoginFields
		if err := readJSON(r.Body, &fields); err != nil {
			resp.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, resp)
			return
		}

		result, err := actions.Login(r.Context(), p, fields)
		if err != nil {
			status := http.StatusBadRequest
			if err == papers.ErrUserNotFound || err == papers.ErrPasswordMismatch {
				resp.Error = "User not found or password doesn't match"
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

		if result.NeedsTOTP {
			if err := p.Config.Storage.Session.MultiSet(r, map[string]string{
				"uid": strconv.FormatInt(result.User.GetID(), 10),
				"rm":  strconv.FormatBool(result.Remember),
			}); err != nil {
				status = http.StatusInternalServerError
				resp.Error = papers.ErrSessionError.Error()
			} else if err := p.Config.Storage.Session.Write(r, w); err != nil {
				status = http.StatusInternalServerError
				resp.Error = papers.ErrSessionError.Error()
			} else {
				resp.NeedsTOTP = true
			}
		} else {
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
			}
		}

		w.WriteHeader(status)
		writeJSON(w, resp)
	}
}
