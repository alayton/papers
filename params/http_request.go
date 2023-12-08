package params

import (
	"net/http"
)

type HTTPRequestParams struct {
}

func (p *HTTPRequestParams) Get(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
