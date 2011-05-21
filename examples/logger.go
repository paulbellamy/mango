package main

import (
	"mango"
)

func Hello(env Env) (Status, Headers, Body) {
	return 200, map[string]string{}, Body("Hello World!")
}

func main() {
	stack := new(mango.Stack)
	stack.address = ":3000"
	stack.Middleware(Logger)
	stack.Run(Hello)
}
