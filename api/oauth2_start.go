package api

import (
	"net/http"
	"strconv"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

type oauth2StartResponse struct {
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

func OAuth2Start(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp oauth2StartResponse

		var fields actions.OAuth2StartFields
		if err := readJSON(r.Body, &fields); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp.Error = err.Error()
			writeJSON(w, resp)
			return
		}

		result, err := actions.OAuth2Start(r.Context(), p, fields)
		if err != nil {
			resp.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, resp)
			return
		}

		status := http.StatusOK
		if err := p.Config.Storage.Session.MultiSet(r, map[string]string{"rm": strconv.FormatBool(fields.Remember), "nonce": result.Nonce}); err != nil {
			status = http.StatusInternalServerError
			resp.Error = papers.ErrSessionError.Error()
		} else {
			resp.URL = result.RedirectURL
		}

		w.WriteHeader(status)
		writeJSON(w, resp)
	}
}
