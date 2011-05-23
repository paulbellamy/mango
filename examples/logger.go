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
	custom_logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	stack.Middleware(mango.Logger(&custom_logger))
	stack.Run(Hello)
}
