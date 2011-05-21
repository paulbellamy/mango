package main

import (
  "mango"
)

func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  return 200, map[string]string{}, mango.Body("Hello World!")
}

func main() {
  app := new(mango.Mango)
  app.Address = ":3000"
  app.Run(Hello)
}
