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

func methodRouter(method string) func(string, http.HandlerFunc) http.HandlerFunc {
	return func(route string, app http.HandlerFunc) http.HandlerFunc {
		regex := regexp.MustCompile(route)
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				return
			}

			matches := regex.FindStringSubmatch(r.URL.Path)
			if len(matches) == 0 {
				return
			}

			r.Header["Route-Matches"] = matches
			app(w, r)
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

var ANY = func(route string, app http.HandlerFunc) http.HandlerFunc {
	regex := regexp.MustCompile(route)
	return func(w http.ResponseWriter, r *http.Request) {
		matches := regex.FindStringSubmatch(r.URL.Path)
		if len(matches) != 0 {
			r.Header["Route-Matches"] = matches
			app(w, r)
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
