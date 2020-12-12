package helpers

import (
	"os"
	"strings"
)

var GetIPFSGateway = func() string {
	return os.Getenv("IPFS_GATEWAY")
}

var DeterminePath = func(path string) string {
	if path == "/" {
		return "/index.html"
	}

	return path
}

var CheckTrailingSpace = func(path string) bool {
	return strings.HasSuffix(path, "/")
}

var TrimRightPath = func(path string) string {
	return strings.TrimRight(path, "/")
}

var AddSlashEachString = func(s ...string) string {
	var str strings.Builder

	for _, word := range s {
		str.WriteString("/" + word)
	}

	return str.String()
}

var GetBytes = func(s string) (b []byte) {
	return []byte(s)
}
