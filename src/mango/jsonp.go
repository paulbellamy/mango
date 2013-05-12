package mango

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var jsonp_valid_callback_matcher *regexp.Regexp = regexp.MustCompile("^[a-zA-Z_$][a-zA-Z_0-9$]*([.]?[a-zA-Z_$][a-zA-Z_0-9$]*)*$")

func JSONP(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		callback := r.FormValue("callback")

		if callback != "" {
			if !jsonp_valid_callback_matcher.MatchString(callback) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Content-Length", "11")
				w.WriteHeader(400)
				w.Write([]byte("Bad Request"))
				return
			}
		}

		wrapped := NewBufferedResponseWriter(nil)
		f(wrapped, r)

		status := wrapped.Status
		headers := wrapped.Header()
		body := wrapped.Body
		if callback != "" && strings.Contains(headers.Get("Content-Type"), "application/json") {
			headers.Set("Content-Type", strings.Replace(headers.Get("Content-Type"), "json", "javascript", -1))
			if headers.Get("Content-Length") != "" {
				headers.Set("Content-Length", fmt.Sprintf("%d", len(body.Bytes())+len(callback)+2))
			}
			for key, values := range headers {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(status)
			w.Write([]byte(fmt.Sprintf("%s(", callback)))
			w.Write(body.Bytes())
			w.Write([]byte(")"))
		} else {
			for key, values := range headers {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(status)
			w.Write(body.Bytes())
		}
	}
}
