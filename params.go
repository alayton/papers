package papers

import (
	"net/http"
)

type RouteParams interface {
	Get(r *http.Request, key string) string
}
