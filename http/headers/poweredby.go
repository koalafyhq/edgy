package headers

import "net/http"

// PoweredBy is
func PoweredBy(server string, via string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("server", server)
		w.Header().Set("x-edgy-via", via)

		next.ServeHTTP(w, r)
	})
}
