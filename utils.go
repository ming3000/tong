package tong

import "net/http"

// fix the input path
func fixPath(path string) string {
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return path
}

// GetPath returns RawPath。
// if it's empty returns Path from URL
func parsePath(r *http.Request) string {
	path := r.URL.RawPath
	if path == "" {
		path = r.URL.Path
	}
	return path
}
