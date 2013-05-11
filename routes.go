package mango

import (
	"net/http"
	"regexp"
)

type instrumentedResponseWriter struct {
	wrapped http.ResponseWriter
	Sent    bool
}

func newInstrumentedResponseWriter(w http.ResponseWriter) *instrumentedResponseWriter {
	return &instrumentedResponseWriter{wrapped: w}
}

func (w *instrumentedResponseWriter) Header() http.Header {
	return w.wrapped.Header()
}

func (w *instrumentedResponseWriter) Write(p []byte) (int, error) {
	w.Sent = true
	return w.wrapped.Write(p)
}

func (w *instrumentedResponseWriter) WriteHeader(status int) {
	w.Sent = true
	w.wrapped.WriteHeader(status)
}

func methodRouter(method string) func(route string, f http.HandlerFunc) http.HandlerFunc {
	return func(route string, f http.HandlerFunc) http.HandlerFunc {
		regex := regexp.MustCompile(route)
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == method && regex.MatchString(r.URL.Path) {
				f(w, r)
			}
		}
	}
}

var GET = methodRouter("GET")
var POST = methodRouter("POST")
var PUT = methodRouter("PUT")
var DELETE = methodRouter("DELETE")
var HEAD = methodRouter("HEAD")
var OPTIONS = methodRouter("OPTIONS")
var PATCH = methodRouter("PATCH")

var ANY = func(route string, f http.HandlerFunc) http.HandlerFunc {
	regex := regexp.MustCompile(route)
	return func(w http.ResponseWriter, r *http.Request) {
		if regex.MatchString(r.URL.Path) {
			f(w, r)
		}
	}
}

func Routes(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseWriter := newInstrumentedResponseWriter(w)
		for _, handler := range handlers {
			// compile the matchers
			handler(responseWriter, r)
			if responseWriter.Sent {
				return
			}
		}
	}
}
