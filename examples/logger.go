package main

import (
	"mango"
)

func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
	return 200, mango.Headers{}, mango.Body("Hello World!")
}

func main() {
	stack := new(mango.Stack)
	stack.address = ":3000"
	custom_logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	stack.Middleware(mango.Logger(custom_logger))
	stack.Run(Hello)
}
