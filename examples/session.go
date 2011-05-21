package main

import (
	"mango"
)

func Hello(env Env) (Status, Headers, Body) {
  env.Session()['old_session_attribute'] = nil
  env.Session()['new_session_attribute'] = 'Never Gonna Give You Up'
	return 200, map[string]string{}, Body("Hello World!")
}

func main() {
	stack := new(mango.Stack)
	stack.address = ":3000"
	stack.Middleware(Sessions("my_secret", "my_key"))
	stack.Run(Hello)
}
