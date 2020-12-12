package proxy

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/koalafy/edgy/http/encoding"
	"github.com/koalafy/edgy/internal/helpers"
)

func (transport *Transporter) responseNotFound(res *http.Response, endpoint string, path string, content []byte) (*http.Response, error) {
	res.StatusCode = 404
	res.Body = ioutil.NopCloser(bytes.NewReader(content))

	return res, nil
}

func (transport *Transporter) responseOK(endpoint string, path string, req *http.Request, res *http.Response) (*http.Response, error) {
	useBrotli := req.Header.Get("x-brotli")
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if useBrotli == "true" {
		buff := new(bytes.Buffer)
		encoded := encoding.CompressToBrotli(body, buff)

		if err := encoded.Close(); err != nil {
			log.Printf("Error compressing content because %v", err)
		} else {
			compressedContent := buff.Bytes()
			res.Header.Set("Content-encoding", "br")
			res.Body = ioutil.NopCloser(bytes.NewReader(compressedContent))
		}
	} else {
		res.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	// FIXME(@faultable): IPFS is return 301, so we force
	// to return 200 to the client
	res.StatusCode = 200

	return res, err
}

func (transport *Transporter) responseFromOrigin(path string, req *http.Request) (*http.Response, error) {
	endpoint := GetEndpointCtx(req.Context())

	IPFSGateway := helpers.GetIPFSGateway()

	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	req.Host = req.URL.Host
	req.URL.Host = IPFSGateway
	req.URL.Path = helpers.AddSlashEachString("ipfs", endpoint) + helpers.DeterminePath(path)
	req.URL.Scheme = "https"

	res, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	return res, err
}
