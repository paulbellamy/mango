package main

import (
	"../../" // Point this at mango
	"log"
	"os"
)

func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  env.Logger().Println("Got a", env.Request().Method, "request for", env.Request().URL)
	return 200, mango.Headers{}, mango.Body("Hello World!")
}

func main() {
	stack := new(mango.Stack)
	stack.Address = ":3000"
	custom_logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	stack.Middleware(mango.Logger(custom_logger))
	stack.Run(Hello)
}
