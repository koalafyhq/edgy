package headers

import "net/http"

// Cleanup is
func Cleanup(res *http.Response) {
	savedHeader := res.Header

	res.Header = http.Header{}

	res.Header.Set("Content-type", savedHeader.Get("Content-Type"))
	res.Header.Set("x-origin-req-id", savedHeader.Get("x-amz-request-id"))
	res.Header.Set("Etag", savedHeader.Get("Etag"))
}
