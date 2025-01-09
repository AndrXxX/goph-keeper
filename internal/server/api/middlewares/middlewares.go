package middlewares

import "net/http"

type Middleware interface {
	Handler(next http.Handler) http.Handler
}
