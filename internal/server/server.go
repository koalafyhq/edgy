package server

import (
	"net/http"

	"github.com/koalafy/edgy/internal/proxy"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func New() *http.ServeMux {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy := proxy.New()

		proxy.ServeHTTP(w, r)
	})

	router := http.DefaultServeMux

	router.HandleFunc("/__healthcheckz", healthcheck)
	router.Handle("/", h)

	return router
}

func Run(app *http.ServeMux) error {
	return http.ListenAndServe(":3000", app)
}
