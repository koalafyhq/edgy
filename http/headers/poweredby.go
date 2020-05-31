package headers

import "net/http"

// PoweredBy is
func PoweredBy(server string, via string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("server", server)
			w.Header().Set("x-edgy-via", via)

			h.ServeHTTP(w, r)
		})
	}
}
