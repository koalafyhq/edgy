package proxy

import (
	"io/ioutil"
	"net/http"

	"github.com/koalafy/edgy/http/headers"
	"github.com/koalafy/edgy/internal/helpers"
	"github.com/rs/zerolog/log"
)

// Transporter is
type Transporter struct{}

// RoundTrip is
func (transport *Transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	endpoint := GetEndpointCtx(req.Context())
	path := req.URL.Path

	res, err := transport.responseFromOrigin(path, req)

	// probably the request is cancelled
	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode > 400 || res.ContentLength == 0 {
		// maybe this is routed via client-side, try another one?
		// TODO(@faultable): handle what if this /favicon.ico for the shake of efficiency?
		newpath := path + ".html"
		res, err = transport.responseFromOrigin(newpath, req)

		if err != nil {
			log.Error().Err(err).Msg("")
		}

		defer res.Body.Close()

		if res.StatusCode > 400 || res.ContentLength == 0 {
			// probably this is a subpath! try another one, again?
			// TODO(@faultable): handle trailing slash here
			subpath := helpers.TrimRightPath(path)

			if path != "" {
				subpath = path + "/index.html"
			}

			res, err = transport.responseFromOrigin(subpath, req)

			if err != nil {
				log.Error().Err(err).Msg("")
			}

			defer res.Body.Close()

			if res.StatusCode > 400 || res.ContentLength == 0 {
				// ok we're give up, try to request custom 404 page (from framework) instead?
				// TODO(@faultable): handle non-Next.js framework here (dynamically)!
				notFoundPath := "/404.html"
				res, err = transport.responseFromOrigin(notFoundPath, req)

				if err != nil {
					log.Error().Err(err).Msg("")
				}

				defer res.Body.Close()

				custom404NotFound, err := ioutil.ReadAll(res.Body)

				if err != nil {
					log.Error().Err(err).Msg("")
				}

				headers.Cleanup(res)
				return transport.responseNotFound(res, endpoint, path, custom404NotFound)
			}

			headers.Cleanup(res)
			return transport.responseOK(endpoint, path, req, res)
		}

		headers.Cleanup(res)
		return transport.responseOK(endpoint, path, req, res)
	}

	headers.Cleanup(res)
	return transport.responseOK(endpoint, path, req, res)
}
