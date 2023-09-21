package rest

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
			p.Logger.Print("Error getting access token during logout:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var refreshToken *papers.RefreshToken
		if err := p.Config.Storage.Cookies.Read(p.Config.RefreshCookieName, r, &refreshToken); err != nil && err != papers.ErrCookieNotFound {
			p.Logger.Print("Error getting refresh token during logout:", err)
		}

		p.Config.Storage.Cookies.Remove(p.Config.AccessCookieName, w)
		p.Config.Storage.Cookies.Remove(p.Config.RefreshCookieName, w)

		if err := actions.Logout(r.Context(), p, user.GetID(), accessToken, refreshToken); err != nil {
			p.Logger.Print("Error invalidating tokens for logout:", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}
