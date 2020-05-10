package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/koalafy/edgy/http/headers"
	"github.com/koalafy/edgy/internal/cachemanager"
	"github.com/koalafy/edgy/internal/helpers"
	"github.com/koalafy/edgy/internal/logger"
)

// Transporter is
type Transporter struct {
	cachemanager *cachemanager.CacheManager
	mu           sync.RWMutex
}

// RoundTrip is
func (transport *Transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	endpoint := GetEndpointCtx(req.Context())
	path := req.URL.Path

	// we use reqID to track failed request
	reqID, _ := headers.IDFromRequest(req)

	cached, content := transport.checkCachedResponse(endpoint, path)

	if cached {
		logger.Cache("HIT", endpoint, path)

		body := bytes.NewBuffer(content)

		return transport.responseFromCache(body, req)
	}

	start := time.Now()
	res, err := transport.responseFromOrigin(path, req)
	elapsed := time.Since(start)

	// probably the request is cancelled
	if err != nil {
		logger.Error(reqID.String(), err, "proxy:"+path)

		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode > 400 || res.ContentLength == 0 {
		// maybe this is routed via client-side, try another one?
		// TODO(@faultable): handle what if this /favicon.ico for the shake of efficiency?
		newpath := path + ".html"
		res, err = transport.responseFromOrigin(newpath, req)

		if err != nil {
			logger.Error(reqID.String(), err, "proxy:"+path)
		}

		defer res.Body.Close()

		if res.StatusCode > 400 || res.ContentLength == 0 {
			// probably this is a subpath! try another one, again?
			// TODO(@faultable): handle trailing slash here
			subpath := helpers.TrimRightPath(path)
			subpath = path + "/index.html"

			res, err = transport.responseFromOrigin(subpath, req)

			if err != nil {
				logger.Error(reqID.String(), err, "proxy:"+path)
			}

			defer res.Body.Close()

			if res.StatusCode > 400 || res.ContentLength == 0 {
				// ok we're give up, try to request custom 404 page (from framework) instead?
				// TODO(@faultable): handle non-Next.js framework here (dynamically)!
				notFoundPath := "/404.html"
				res, err = transport.responseFromOrigin(notFoundPath, req)

				if err != nil {
					logger.Error(reqID.String(), err, "proxy:"+path)
				}

				defer res.Body.Close()

				custom404NotFound, err := ioutil.ReadAll(res.Body)

				if err != nil {
					logger.Error(reqID.String(), err, "proxy:"+path)
				}

				logger.CacheWithLatency("MISS", endpoint, path, elapsed)

				return transport.responseNotFound(res, endpoint, path, custom404NotFound)
			}
		}
	}

	logger.CacheWithLatency("MISS", endpoint, path, elapsed)

	return transport.responseOK(endpoint, path, res)
}
