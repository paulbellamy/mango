package main

import (
  "../../" // Point this to mango
)

// Our default handler
func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
	return 200, mango.Headers{}, mango.Body("Hello World!")
}

// Our handler for /goodbye
func Goodbye(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
	return 200, mango.Headers{}, mango.Body("Goodbye World!")
}

func main() {
	stack := new(mango.Stack)
	stack.Address = ":3000"

	// Route all requests for /goodbye to the Goodbye handler
	routes := map[string]mango.App{"/goodbye(.*)": Goodbye}
	stack.Middleware(mango.Routing(routes))

	// Hello handles all requests not sent to Goodbye
	stack.Run(Hello)
}
