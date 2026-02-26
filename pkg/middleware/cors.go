package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/puriice/httplibs/pkg/middleware/cors"
)

type config struct {
	origins  map[string]struct{}
	allowAll bool
	headers  string
	methods  string
	expose   string
	maxAge   string
	time     string
}

func Cors(option cors.CorsOptions) Middleware {
	config := &config{
		origins:  make(map[string]struct{}, len(option.AllowOrigins)),
		allowAll: false,
	}

	for _, o := range option.AllowOrigins {
		if o == "*" {
			config.allowAll = true
			continue
		}
		config.origins[o] = struct{}{}
	}

	// Precompute joined strings once
	if len(option.AllowMethods) > 0 {
		config.methods = strings.Join(option.AllowMethods, ", ")
	}

	if len(option.AllowHeaders) > 0 {
		config.headers = strings.Join(option.AllowHeaders, ", ")
	}

	if len(option.AllowExposeHeaders) > 0 {
		config.expose = strings.Join(option.AllowExposeHeaders, ", ")
	}

	if len(option.TimingAllowOrigin) > 0 {
		config.time = strings.Join(option.TimingAllowOrigin, ", ")
	}

	if option.MaxAge > 0 {
		config.maxAge = strconv.Itoa(option.MaxAge)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("Cross-Origin-Resource-Policy", "cross-origin")

			origin := r.Header.Get("Origin")

			if origin == "" {
				if option.AllowNoOrigin {
					next.ServeHTTP(w, r)
					return
				}

				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !config.allowAll {
				if _, ok := config.origins[origin]; !ok {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}

			h.Add("Vary", "Origin")
			h.Set("Access-Control-Allow-Origin", origin)

			if option.AllowCredentials {
				h.Set("Access-Control-Allow-Credentials", "true")
			}

			if len(option.AllowExposeHeaders) > 0 {
				h.Set("Access-Control-Expose-Headers", config.expose)
			}

			if option.MaxAge > 0 {
				h.Set("Access-Control-Max-Age", config.maxAge)
			}

			if len(option.TimingAllowOrigin) > 0 {
				h.Set("Timing-Allow-Origin", config.time)
			}

			if r.Method != http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			if header := r.Header.Get("Access-Control-Request-Headers"); header != "" {
				if len(option.AllowHeaders) > 0 {
					h.Set("Access-Control-Allow-Headers", config.headers)
				} else {
					h.Set("Access-Control-Allow-Headers", header)
				}
			}

			if option.AllowMethods != nil {
				h.Set("Access-Control-Allow-Methods", config.methods)
			} else {
				h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			}

			w.WriteHeader(http.StatusNoContent)
		})
	}
}
