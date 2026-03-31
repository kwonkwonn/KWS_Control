package api

import "net/http"

// withSecurityHeaders wraps a handler to add common security headers to every response.
func withSecurityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next(w, r)
	}
}
