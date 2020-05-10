package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/koalafy/edgy/http/headers"
	"github.com/koalafy/edgy/internal/cachemanager"
	"github.com/koalafy/edgy/internal/helpers"
	"github.com/koalafy/edgy/internal/logger"
)

// Router is
type Router struct {
	Proxy        *httputil.ReverseProxy
	CacheManager *cachemanager.CacheManager
}

// New is
func New(cacheManager *cachemanager.CacheManager) *Router {
	modifyResponse := func(res *http.Response) (err error) {
		return nil
	}

	transport := &Transporter{
		cachemanager: cacheManager,
	}

	router := &Router{
		CacheManager: cacheManager,
		Proxy: &httputil.ReverseProxy{
			Transport:      transport,
			ModifyResponse: modifyResponse,
			Director:       func(req *http.Request) {},
		},
	}

	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	host := r.Host

	// we use reqID to track failed request
	reqID, _ := headers.IDFromRequest(r)

	endpoint, err := router.CacheManager.GetEndpoint(host)

	// does current request host (endpoint) exist on our db?
	// if not, return not found
	if err != nil {
		http.NotFound(w, r)
		logger.Error(reqID.String(), err, host)

		return
	}

	endpointCtx := SetEndpointCtx(r.Context(), string(endpoint))

	if r.Method == "GET" {
		// check trailing slash
		l := len(path) - 1

		// if any, remove the trailing slash
		// by redirect it
		if l > 0 && helpers.CheckTrailingSpace(path) {
			path = path[:l]
			uri := path

			// TODO(@faultable): handle request cancellation
			http.Redirect(w, r, uri, http.StatusMovedPermanently)

			return
		}

		router.Proxy.ServeHTTP(w, r.WithContext(endpointCtx))

		return
	}
}
