package helpers

import "strings"

// DeterminePath is
var DeterminePath = func(path string) string {
	if path == "/" {
		return "/index.html"
	}

	return path
}

// CheckTrailingSpace is
var CheckTrailingSpace = func(path string) bool {
	return strings.HasSuffix(path, "/")
}

// TrimRightPath is
var TrimRightPath = func(path string) string {
	return strings.TrimRight(path, "/")
}

// AddSlashEachString is helper to add slash into each word
var AddSlashEachString = func(s ...string) string {
	var str strings.Builder

	for _, word := range s {
		str.WriteString("/" + word)
	}

	return str.String()
}

// GetBytes is a helper to convert `string` into `byte`
var GetBytes = func(s string) (b []byte) {
	return []byte(s)
}
