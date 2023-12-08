package chi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ChiRouteParams struct {
}

func (p *ChiRouteParams) Get(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
