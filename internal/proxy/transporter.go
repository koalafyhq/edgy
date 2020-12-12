package proxy

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/koalafy/edgy/internal/helpers"
)

type Transporter struct{}

func (transport *Transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	endpoint := GetEndpointCtx(req.Context())
	path := req.URL.Path

	res, err := transport.responseFromOrigin(path, req)

	// probably the request is cancelled
	if err != nil {
		log.Println(err)

		return nil, err
	}

	defer res.Body.Close()

	// FIXME(@faultable)?: This is IPFS specific logic
	// to determine the most CORRECT requested path
	// probably myself in future should fix this?
	//
	// Myself in future (12/12/20): NO, I have no idea.
	//
	getEtag := res.Header.Get("Etag")
	isDirListing := "\"DirIndex"
	checkEtag := ""

	if len(getEtag) > 8 {
		checkEtag = getEtag[:9]
	}

	// This is a stupid hack like nginx try_files $uri $uri/ $uri.html
	// but for reverse proxy and in stupid way
	if res.StatusCode > 400 || checkEtag == isDirListing {
		// maybe this is (previously) routed via client-side, try another one?
		// e.g: request: /about, upstream: /about.html
		newpath := path + ".html"
		res, err = transport.responseFromOrigin(newpath, req)

		if err != nil {
			log.Println(err)
		}

		defer res.Body.Close()

		if res.StatusCode > 400 || res.ContentLength == 0 {
			// probably this is a subpath! try another one, again?
			subpath := helpers.TrimRightPath(path)
			// e.g: request: /about, upstream: /about/index.html
			if path != "" {
				subpath = path + "/index.html"
			}

			res, err = transport.responseFromOrigin(subpath, req)

			if err != nil {
				log.Println(err)
			}

			defer res.Body.Close()

			if res.StatusCode > 400 || res.ContentLength == 0 {
				// ok we're give up, try to request custom 404 page (from framework) instead?
				// TODO(@faultable): handle non-Next.js framework here (dynamically)!
				notFoundPath := "/404.html"
				res, err = transport.responseFromOrigin(notFoundPath, req)

				if err != nil {
					log.Println(err)
				}

				defer res.Body.Close()

				custom404NotFound, err := ioutil.ReadAll(res.Body)

				if err != nil {
					log.Println(err)
				}

				return transport.responseNotFound(res, endpoint, path, custom404NotFound)
			}

			return transport.responseOK(endpoint, path, req, res)
		}

		return transport.responseOK(endpoint, path, req, res)
	}

	return transport.responseOK(endpoint, path, req, res)
}
