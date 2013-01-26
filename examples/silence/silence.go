package main

import (
  "../../" // Point this to mango
)

// Our custom middleware
func SilenceErrors(env mango.Env, app mango.App) (mango.Status, mango.Headers, mango.Body) {
	// Call our upstream app
	status, headers, body := app(env)

	// If we got an error
	if status == 500 {
		// Silence it!
		status = 200
		headers = mango.Headers{}
		body = "Silence is golden!"
	}

	// Pass the response back to the client
	return status, headers, body
}

// Our default handler
func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  //Return 500 to trigger the silence
	return 500, mango.Headers{}, mango.Body("Hello World!")
}

func main() {
  stack := new(mango.Stack)
  stack.Address = ":3000"

  stack.Middleware(SilenceErrors) // Include our custom middleware

  stack.Run(Hello)
}
