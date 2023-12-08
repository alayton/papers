package api

import (
	"errors"
	"net/http"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

type recoverResponse struct {
	Error string `json:"error,omitempty"`
}

func RecoverStart(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp recoverResponse

		var fields actions.RecoverStartFields
		if err := readJSON(r.Body, &fields); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp.Error = err.Error()
			writeJSON(w, resp)
			return
		}

		err := actions.Recover(r.Context(), p, fields)
		if err != nil {
			p.Logger.Error("recovery action errored", "error", err)
			if err == papers.ErrMissingEmail || err == papers.ErrInvalidEmail {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			unwrapped := errors.Unwrap(err)
			if unwrapped != nil {
				resp.Error = unwrapped.Error()
			} else {
				resp.Error = err.Error()
			}

			writeJSON(w, resp)
			return
		}

		writeJSON(w, resp)
	}
}
