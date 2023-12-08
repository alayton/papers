package api

import (
	"net/http"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

func Confirm(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp recoverResponse

		var fields actions.ConfirmFields
		fields.Token = p.Config.RouteParams.Get(r, "token")

		err := actions.Confirm(r.Context(), p, fields)
		if err == papers.ErrTokenNotFound {
			w.WriteHeader(http.StatusBadRequest)
			resp.Error = err.Error()
		} else if err != nil {
			p.Logger.Error("confirm action errored", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Error = papers.ErrStorageError.Error()
		}

		writeJSON(w, resp)
	}
}
