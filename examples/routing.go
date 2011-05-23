package main

import (
	"mango"
)

// Our default handler
func Hello(env Env) (Status, Headers, Body) {
	return 200, map[string]string{}, Body("Hello World!")
}

// Our handler for /goodbye
func Goodbye(env Env) (Status, Headers, Body) {
	return 200, map[string]string{}, Body("Goodbye World!")
}

func main() {
	stack := new(mango.Stack)
	stack.address = ":3000"

	// Route all requests for /goodbye to the Goodbye handler
	routes := map[string]mango.App{"/goodbye(.*)": Goodbye}
	stack.Middleware(Routing(routes))

	// Hello handles all requests not sent to Goodbye
	stack.Run(Hello)
}
