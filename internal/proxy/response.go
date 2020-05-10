package proxy

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/koalafy/edgy/http/headers"
	"github.com/koalafy/edgy/internal/helpers"
	"github.com/koalafy/edgy/internal/logger"
)

func (transport *Transporter) checkCachedResponse(endpoint string, path string) (bool, []byte) {
	cacheKey := endpoint + path

	transport.mu.RLock()
	defer transport.mu.RUnlock()

	cached, err := transport.cachemanager.Get(cacheKey)

	if err != nil {
		return false, nil
	}

	return true, cached
}

func (transport *Transporter) responseNotFound(res *http.Response, endpoint string, path string, content []byte) (*http.Response, error) {
	cacheKey := endpoint + path

	transport.mu.RLock()
	defer transport.mu.RUnlock()

	headers.Cleanup(res)

	// add our header to track the cache ratio easier
	// this header will be stored on cache
	res.Header.Set("x-edgy-cache", "HIT")

	res.StatusCode = 404
	res.Body = ioutil.NopCloser(bytes.NewReader(content))

	body, err := httputil.DumpResponse(res, true)

	if err != nil {
		logger.Error("", err, "proxy:DumpResponse")
	}

	if res.ContentLength < 10000000 && res.ContentLength > 0 {
		transport.cachemanager.Set(cacheKey, body)

		logger.CacheWithSize("ADDED", endpoint, path, res.ContentLength)
	}

	// add our header to track the cache ratio easier
	// this header will be sent once to the client
	res.Header.Set("x-edgy-cache", "MISS")

	return res, nil
}

func (transport *Transporter) responseOK(endpoint string, path string, res *http.Response) (*http.Response, error) {
	cacheKey := endpoint + path

	transport.mu.Lock()
	defer transport.mu.Unlock()

	body, err := httputil.DumpResponse(res, true)

	if err != nil {
		logger.Error("", err, "proxy:DumpResponse")
	}

	// We don't want to store to the cache if payload > 10MB
	// And we don't want to store empty response, right?
	if res.ContentLength > 0 && res.ContentLength < 10000000 {
		transport.cachemanager.Set(cacheKey, body)

		logger.CacheWithSize("ADDED", endpoint, path, res.ContentLength)
	}

	// add our header to track the cache ratio easier
	// this header will be sent once to the client
	res.Header.Set("x-edgy-cache", "MISS")

	return res, err
}

func (transport *Transporter) responseFromCache(buff *bytes.Buffer, request *http.Request) (*http.Response, error) {
	return http.ReadResponse(bufio.NewReader(buff), request)
}

func (transport *Transporter) responseFromOrigin(path string, req *http.Request) (*http.Response, error) {
	endpoint := GetEndpointCtx(req.Context())

	// TODO(@faultable): fix this
	s3Gateway := os.Getenv("S3_GATEWAY")

	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	req.Host = req.URL.Host
	req.URL.Host = s3Gateway
	req.URL.Path = helpers.AddSlashEachString("edgy-bundles", endpoint) + helpers.DeterminePath(path)
	req.URL.Scheme = "https"

	res, err := http.DefaultTransport.RoundTrip(req)

	// remove unecessary headers
	headers.Cleanup(res)

	// add our header to track the cache ratio easier
	// this header will be stored on cache
	res.Header.Set("x-edgy-cache", "HIT")

	return res, err
}
