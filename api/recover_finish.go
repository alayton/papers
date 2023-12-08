package api

import (
	"errors"
	"net/http"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

func RecoverFinish(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp recoverResponse

		var fields actions.RecoverFinishFields
		if err := readJSON(r.Body, &fields); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp.Error = err.Error()
			writeJSON(w, resp)
			return
		}

		err := actions.RecoverFinish(r.Context(), p, fields)
		if err == papers.ErrTokenNotFound || errors.Is(err, papers.ErrInvalidPassword) {
			w.WriteHeader(http.StatusBadRequest)
			resp.Error = err.Error()
		} else if err != nil {
			p.Logger.Error("recovery action errored", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = papers.ErrStorageError.Error()
		}

		writeJSON(w, resp)
	}
}
