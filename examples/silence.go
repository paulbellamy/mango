package main

import (
  "mango"
  "silence_middleware"
)

// Our default handler
func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  //Return 500 to trigger the silence
	return 500, mango.Headers{}, mango.Body("Hello World!")
}

func main() {
  stack := new(mango.Stack)
  stack.Address = ":3000"

  stack.Middleware(silence.SilenceErrors) // Include our custom middleware

  stack.Run(Hello)
}