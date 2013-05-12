// Mango is a modular web-application framework for Go, inspired by Rack and PEP333.
package mango

import (
	"bytes"
	"net/http"
)

type BufferedResponseWriter struct {
	wrapped http.ResponseWriter
	headers http.Header
	started bool
	Status  int
	Body    bytes.Buffer
}

func NewBufferedResponseWriter(w http.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{
		headers: make(map[string][]string),
		wrapped: w,
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

// Send data to the client
func (w *BufferedResponseWriter) Flush() {
	if w.started {
		output_headers := w.wrapped.Header()
		for name, values := range w.headers {
			for _, value := range values {
				output_headers.Add(name, value)
			}
		}

		w.wrapped.WriteHeader(w.Status)
		w.wrapped.Write(w.Body.Bytes())
	}
}
