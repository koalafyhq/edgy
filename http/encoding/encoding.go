package encoding

import (
	"bytes"

	"github.com/andybalholm/brotli"
)

// CompressToBrotli —
func CompressToBrotli(content []byte, buff *bytes.Buffer) *brotli.Writer {
	compressed := brotli.NewWriter(buff)

	compressed.Write(content)

	return compressed
}
