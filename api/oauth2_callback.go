package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

type oauth2CallbackResponse struct {
	User  papers.User `json:"user,omitempty"`
	Error string      `json:"error,omitempty"`
}

func OAuth2Callback(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp oauth2CallbackResponse

		nonce, err := p.Config.Storage.Session.Get(r, "nonce")
		if err != nil {
			resp.Error = papers.ErrStorageError.Error()
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, resp)
			return
		}
		rm, _ := p.Config.Storage.Session.Get(r, "rm")

		p.Config.Storage.Session.MultiDelete(r, []string{"nonce", "rm"})

		didError := r.FormValue("error")
		if len(didError) > 0 {
			reason := r.FormValue("error_reason")
			p.Logger.Error("oauth2 login failed", "error", reason)
			resp.Error = papers.ErrOAuth2LoginFailed.Error()
			w.WriteHeader(http.StatusUnauthorized)
			writeJSON(w, resp)
			return
		}

		fields := actions.OAuth2CallbackFields{
			Provider: strings.ToLower(filepath.Base(r.URL.Path)),
			State:    r.FormValue("state"),
			Code:     r.FormValue("code"),
			Nonce:    nonce,
			Remember: rm == "true",
		}

		result, err := actions.OAuth2Callback(r.Context(), p, fields)
		if err != nil {
			resp.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
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
		}

		w.WriteHeader(status)
		writeJSON(w, resp)
	}
}
