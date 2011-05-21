package main

import (
  "mango"
)

func Hello(env Env) (Status, Headers, Body) {
	return 200, map[string]string{}, Body("Hello World!")
}

func main() {
	mango := new(Mango)
	mango.address = ":3000"
  mango.Middleware(Logger)
	mango.Run(Hello)
}
