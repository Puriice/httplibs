package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Pipe(middleware Middleware, middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for i := len(middlewares) - 1; i >= 0; i-- {
				next = middlewares[i](next)
			}

			next = middleware(next)

			next.ServeHTTP(w, r)
		})
	}
}
