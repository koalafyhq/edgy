package server

import (
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
	"github.com/koalafy/edgy/http/headers"
	"github.com/koalafy/edgy/internal/cachemanager"
	"github.com/koalafy/edgy/internal/proxy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// New is
func New(cacheMananger *cachemanager.CacheManager) *http.ServeMux {
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	via := os.Getenv("EDGY_REGION")
	reqIDKey := "x-edgy-req-id"
	serverName := "edgy"

	c := alice.New()
	c = c.Append(hlog.NewHandler(log))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("host", r.Host).
			Str("via", via).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	c = c.Append(hlog.RequestIDHandler("req_id", reqIDKey))
	c = c.Append(headers.PoweredBy(serverName, via))

	h := c.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy := proxy.New(cacheMananger)

		proxy.ServeHTTP(w, r)
	}))

	router := http.DefaultServeMux

	router.HandleFunc("/__healthcheckz", healthcheck)
	router.Handle("/", h)

	return router
}

// Run is
func Run(app *http.ServeMux) error {
	log.Debug().Msg("Server booted")

	return http.ListenAndServe(":3000", app)
}
