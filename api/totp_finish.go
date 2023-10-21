package api

import (
	"errors"
	"net/http"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

type totpFinishResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func TOTPFinish(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, loggedIn := p.LoggedInUser(r)
		if !loggedIn {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var resp totpFinishResponse

		var fields actions.TOTPFinishFields
		if err := readJSON(r.Body, &fields); err != nil {
			resp.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, resp)
			return
		}

		err := actions.TOTPFinish(r.Context(), p, user, fields)
		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, papers.ErrStorageError) || errors.Is(err, papers.ErrCryptoError) {
				status = http.StatusInternalServerError
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

		resp.Success = true

		writeJSON(w, resp)
	}
}
