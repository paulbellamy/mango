package main

import (
	. "mango"
	"net/http"
)

// Our custom middleware
func SilenceErrors(app http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// catch the response
		response := NewBufferedResponseWriter(w)

		// Call our upstream app
		app(response, r)

		// If we got an error
		if response.Status == 500 {
			// Silence it!
			response.Status = 200
			response.Body.Reset()
			response.Body.Write([]byte("Silence is golden!"))
		}

		// Send our output to the client
		response.Flush()
	}
}

// Our default handler
func Hello(w http.ResponseWriter, r *http.Request) {
	//Return 500 to trigger the silence
	w.WriteHeader(500)
	w.Write([]byte("Hello World!"))
}

func main() {
	http.HandleFunc("/", SilenceErrors(Hello))
	http.ListenAndServe(":3000", nil)
}
