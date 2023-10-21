package api

import (
	"net/http"

	"github.com/alayton/papers"
	"github.com/alayton/papers/actions"
)

type registerResponse struct {
	User   papers.User `json:"user,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

func Register(p *papers.Papers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp registerResponse

		var fields actions.RegisterFields
		if err := readJSON(r.Body, &fields); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp.Errors = []string{err.Error()}
			writeJSON(w, resp)
			return
		}

		user, err := actions.Register(r.Context(), p, fields)
		if err != nil {
			if errs, ok := err.(actions.MultiError); ok {
				for _, e := range errs.Unwrap() {
					resp.Errors = append(resp.Errors, e.Error())
				}
			} else {
				resp.Errors = []string{err.Error()}
			}
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, resp)
			return
		}

		resp.User = user
		writeJSON(w, resp)
	}
}
