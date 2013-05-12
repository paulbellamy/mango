package main

import (
	. "../../" // Point this to mango
	"net/http"
)

// Our default handler
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// Our handler for /goodbye
func Goodbye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Goodbye World!"))
}

func main() {
	// Route all requests for /goodbye to the Goodbye handler
	// all other requests go to Hello
	app := Routes(
		GET("/goodbye(.*)", Goodbye),
		ANY(".*", Hello),
	)

	http.HandleFunc("/", app)
	http.ListenAndServe(":3000", nil)
}
