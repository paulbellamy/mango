package silence

import (
	"mango"
)

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
