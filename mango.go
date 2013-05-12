// Mango is a modular web-application framework for Go, inspired by Rack and PEP333.
package mango

import (
	"bytes"
	"fmt"
	"net/http"
)

type BufferedResponseWriter struct {
	headers http.Header
	started bool
	Status  int
	Body    bytes.Buffer
}

func NewBufferedResponseWriter() *BufferedResponseWriter {
	return &BufferedResponseWriter{
		headers: make(map[string][]string),
	}
}

func (w *BufferedResponseWriter) Header() http.Header {
	return w.headers
}

func (w *BufferedResponseWriter) Write(p []byte) (int, error) {
	if !w.started {
		w.started = true
		w.Status = 200
	}
	return w.Body.Write(p)
}

func (w *BufferedResponseWriter) WriteHeader(status int) {
	if !w.started {
		w.started = true
		w.Status = status
	}
}

func init() {
	fmt.Println("Testing Mango Version:", VersionString())
}

func Version() []int {
	return []int{0, 5, 0}
}

func VersionString() string {
	v := Version()
	return fmt.Sprintf("%d.%02d.%02d", v[0], v[1], v[2])
}
