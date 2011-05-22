package main

import (
	"mango"
)

func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
	// to add a session attribute just add it to the map
	env.Session()["new_session_attribute"] = "Never Gonna Give You Up"

	// To remove a session attribute delete it from the map
	env.Session()["old_session_attribute"] = nil, false

	return 200, map[string]string{}, mango.Body("Hello World!")
}

func main() {
	stack := new(mango.Stack)
	stack.Address = ":3000"
	stack.Middleware(mango.Sessions("my_secret", "my_session_key"))
	stack.Run(Hello)
}
