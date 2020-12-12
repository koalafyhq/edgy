package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/koalafy/edgy/internal/helpers"
)

type Router struct {
	Proxy *httputil.ReverseProxy
}

func New() *Router {
	transport := &Transporter{}

	router := &Router{
		Proxy: &httputil.ReverseProxy{
			Transport: transport,
			Director:  func(req *http.Request) {},
		},
	}

	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	endpoint := r.Header.Get("X-IPFS-PATH")
	endpointCtx := SetEndpointCtx(r.Context(), string(endpoint))

	// someone is bypassed the edge!
	if endpoint == "" {
		http.NotFound(w, r)

		return
	}

	// only server GET method, who use another http verb for static site, anyways?
	if r.Method == "GET" {
		// check trailing slash
		l := len(path) - 1

		// if any, remove the trailing slash
		// by redirect it since we never want to access as a directory
		if l > 0 && helpers.CheckTrailingSpace(path) {
			path = path[:l]
			uri := path

			http.Redirect(w, r, uri, http.StatusMovedPermanently)

			return
		}

		router.Proxy.ServeHTTP(w, r.WithContext(endpointCtx))

		return
	}
}
