package server

import (
	"net/http"
	"os"

	"github.com/koalafy/edgy/http/headers"
	"github.com/koalafy/edgy/internal/cachemanager"
	"github.com/koalafy/edgy/internal/logger"
	"github.com/koalafy/edgy/internal/proxy"
)

// New is
func New(cacheMananger *cachemanager.CacheManager) *http.ServeMux {
	reqIDKey := "x-edgy-req-id"
	serverName := "edgy"

	router := http.DefaultServeMux
	proxy := proxy.New(cacheMananger)

	via := os.Getenv("EDGY_REGION")

	router.Handle("/", headers.RequestID(reqIDKey, headers.PoweredBy(serverName, via, proxy)))

	return router
}

// Run is
func Run(app *http.ServeMux) error {
	logger.Debug("Server booted")

	return http.ListenAndServeTLS(":3000", "certs/public.crt", "certs/private.key", app)
}
