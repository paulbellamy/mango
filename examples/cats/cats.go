package main

import (
	. "../../" // Point this to mango
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
)

// Our custom middleware
func Cats(cat_images []string, app http.HandlerFunc) http.HandlerFunc {
	// Initial setup stuff here
	// Done on application setup

	// Initialize our regex for finding image links
	regex := regexp.MustCompile("[^\"']+(.jpg|.png|.gif)")

	// This is our middleware's request handler
	return func(w http.ResponseWriter, r *http.Request) {
		// create a buffer to catch the response
		response := NewBufferedResponseWriter(w)

		// Call the upstream application
		app(response, r)

		// Pick a random cat image
		image_url := cat_images[rand.Int()%len(cat_images)]

		// Substitute in our cat picture. There is probably a more efficient way of
		// doing this, but in the interest of brevity we just replace the whole
		// body, which is inefficient.
		newBody := regex.ReplaceAllString(response.Body.String(), image_url)
		response.Body.Reset()
		response.Body.WriteString(newBody)

		// Send the modified response onwards
		response.Flush()
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Got a", r.Method, "request for", r.URL)

	// Fetch some dummy html
	response, err := http.Get("http://www.example.com/")
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	w.Write(contents)
}

func main() {

	// Initialize our cats middleware with our list of cat_images
	cat_images := []string{"http://images.cheezburger.com/completestore/2010/7/4/9440dc57-52a6-4122-9ab3-efd4daa0ff60.jpg", "http://images.icanhascheezburger.com/completestore/2008/12/10/128733944185267668.jpg"}
	app := Cats(cat_images, Hello)

	http.HandleFunc("/", app)
	http.ListenAndServe(":3000", nil)
}
