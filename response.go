package tong

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// Response wraps the http.ResponseWriter and implements its interface,
// it is used by HTTP handler to generate an HTTP response
type Response struct {
	Writer          http.ResponseWriter
	Status          int
	Size            int
	IfHeaderBeenSet bool
}

// NewResponse create a new instance of Response
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{Writer: w, Status: http.StatusOK, Size: 0, IfHeaderBeenSet: false}
}

// Reset reset the Response instance
func (r *Response) Reset(w http.ResponseWriter) {
	r.Writer = w
	r.Status = http.StatusOK
	r.Size = 0
	r.IfHeaderBeenSet = false
}

// Header returns the http.header map of the writer
func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

// WriteHeader set the HTTP response header with status code
func (r *Response) WriteHeader(code int) {
	if r.IfHeaderBeenSet {
		return
	} // if>

	r.Writer.WriteHeader(code)
	r.Status = code
	r.IfHeaderBeenSet = true
}

// Write writes the data to the client
func (r *Response) Write(data []byte) (int, error) {
	// if Header has not been set,
	// the status will be the default value as StatusOK
	if !r.IfHeaderBeenSet {
		r.WriteHeader(r.Status)
	} // if>

	n, err := r.Writer.Write(data)
	r.Size += n
	return n, err
}

// https://golang.org/pkg/net/http/#Flusher
func (r *Response) Flush() {
	if flusher, ok := r.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

// https://golang.org/pkg/net/http/#Hijacker
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := r.Writer.(http.Hijacker); ok {
		return hijacker.Hijack()
	} else {
		return nil, nil, errors.New("reflect Hijacker error")
	} // else>
}
