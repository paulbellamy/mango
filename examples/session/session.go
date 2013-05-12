package main

import (
	. "../../" // Point this to mango
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	// Load the session from the request
	session := Session(r, "my_session_key", "my_secret", nil)

	// to add a session attribute
	session.Set("new_session_attribute", "Never Gonna Give You Up")

	// to read a session atribute
	counter, _ := session.Get("counter").(int)
	session.Set("counter", counter+1)

	// To remove a session attribute
	session.Del("old_session_attribute")

	// Finish the session before sending a response
	session.Write(w)
	w.Write([]byte(fmt.Sprintf("Session contained: %v", session)))
}

func main() {
	http.HandleFunc("/", Hello)
	http.ListenAndServe(":3000", nil)
}
