package headers

import "net/http"

// Cleanup is
func Cleanup(res *http.Response) {
	res.Header.Del("x-amz-request-id")
	res.Header.Del("server")
}
