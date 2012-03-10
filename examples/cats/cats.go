package main

import (
  "../../" // Point this to mango
  "net/http"
  "os"
  "io/ioutil"
	"math/rand"
	"regexp"
)

// Our custom middleware
func Cats(cat_images []string) mango.Middleware {
	// Initial setup stuff here
	// Done on application setup

	// Initialize our regex for finding image links
	regex := regexp.MustCompile("[^\"']+(.jpg|.png|.gif)")

	// This is our middleware's request handler
	return func(env mango.Env, app mango.App) (mango.Status, mango.Headers, mango.Body) {
		// Call the upstream application
		status, headers, body := app(env)

		// Pick a random cat image
		image_url := cat_images[rand.Int()%len(cat_images)]

		// Substitute in our cat picture
		body = mango.Body(regex.ReplaceAllString(string(body), image_url))

		// Send the modified response onwards
		return status, headers, body
	}
}

func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  env.Logger().Println("Got a", env.Request().Method, "request for", env.Request().URL)

  response, err := http.Get("http://www.example.com/")
  if err != nil {
  	env.Logger().Printf("%s", err)
    os.Exit(1)
  }
  defer response.Body.Close()
  
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    env.Logger().Printf("%s", err)
    os.Exit(1)
  }
  
  return 200, mango.Headers{}, mango.Body(contents)
}

func main() {

  stack := new(mango.Stack)
  stack.Address = ":3000"

  // Initialize our cats middleware with our list of cat_images
  cat_images := []string{"http://images.cheezburger.com/completestore/2010/7/4/9440dc57-52a6-4122-9ab3-efd4daa0ff60.jpg", "http://images.icanhascheezburger.com/completestore/2008/12/10/128733944185267668.jpg"}
  cats_middleware := Cats(cat_images)

  stack.Middleware(cats_middleware) // Include the Cats middleware in our stack

  stack.Run(Hello)
}
