package middleware

import "net/http"

func SecureConnection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add HSTS Header
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}
