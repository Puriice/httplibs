package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/puriice/httplibs/pkg/middleware/cors"
)

func Cors(option cors.CorsOptions) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")

			origin := r.Header.Get("Origin")

			if origin == "" {
				if option.AllowNoOrigin {
					next.ServeHTTP(w, r)
					return
				}

				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !slices.Contains(option.AllowOrigins, origin) && !slices.Contains(option.AllowOrigins, "*") {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if option.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if len(option.AllowExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(option.AllowExposeHeaders, ", "))
			}

			if option.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", fmt.Sprint(option.MaxAge))
			}

			if len(option.TimingAllowOrigin) > 0 {
				w.Header().Set("Timing-Allow-Origin", strings.Join(option.TimingAllowOrigin, ", "))
			}

			if r.Method != http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			if header := r.Header.Get("Access-Control-Request-Headers"); header != "" {
				if len(option.AllowHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(option.AllowHeaders, ", "))
				} else {
					w.Header().Set("Access-Control-Allow-Headers", header)
				}
			}

			if option.AllowMethods != nil {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(option.AllowMethods, ", "))
			} else {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			}

			w.WriteHeader(http.StatusNoContent)
		})
	}
}
