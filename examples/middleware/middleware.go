package main

import (
	. "../../" // Point this to mango
	"log"
	"net/http"
	"time"
)

// Let's define a custom middleware to do some logging
func Logger(app http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		app(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(start))
	}
}

// An upstream app to call
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func main() {

	// Initialize our stack of middleware
  root := http.Dir("./static")
	app :=
		Logger(
			ShowErrors("",
				Static(root,
					Hello)))

	http.HandleFunc("/", app)
	http.ListenAndServe(":3000", nil)
}
