package main

import (
  "mango"
  "cats_middleware"
  "net/http"
  "os"
  "io/ioutil"
)

func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  env.Logger().Println("Got a", env.Request().Method, "request for", env.Request().RawURL)

  response, err := http.Get("http://www.weddinggalleryweb.com/")
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
  cats_middleware := cats.Cats(cat_images)

  stack.Middleware(cats_middleware) // Include the Cats middleware in our stack

  stack.Run(Hello)
}
