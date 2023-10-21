package api

import (
	"net/http"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

func Logout(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, loggedIn := p.LoggedInUser(r)
		if !loggedIn {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var accessToken *papers.AccessToken
		if err := p.Config.Storage.Cookies.Read(p.Config.AccessCookieName, r, &accessToken); err != nil {
			p.Logger.Error("failed to read access token during logout", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var refreshToken *papers.RefreshToken
		if err := p.Config.Storage.Cookies.Read(p.Config.RefreshCookieName, r, &refreshToken); err != nil && err != papers.ErrCookieNotFound {
			p.Logger.Error("failed to read refresh token during logout", "error", err)
		}

		p.Config.Storage.Cookies.Remove(p.Config.AccessCookieName, w)
		p.Config.Storage.Cookies.Remove(p.Config.RefreshCookieName, w)

		if err := actions.Logout(r.Context(), p, user.GetID(), accessToken, refreshToken); err != nil {
			p.Logger.Error("failed to invalidate tokens during logout", "error", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}
