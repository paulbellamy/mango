package cats

import (
	"mango"
	"rand"
	"regexp"
)

func Cats(cat_images []string) mango.Middleware {
	// Initial setup stuff here
	// Done on application setup

	// Initialize our regex for finding image links
	regex := regexp.MustCompile("[\"']*(.jpg|.png)[\"']")

	// This is our middleware's request handler
	return func(env mango.Env, app mango.App) (mango.Status, mango.Headers, mango.Body) {
		// Call the upstream application
		status, headers, body := app(env)

		// Pick a random cat image
		image_url := cat_images[rand.Int()%len(cat_images)]

		// Substitute in our cat picture
		body = mango.Body(regex.ReplaceAllString(string(body), image_url))

		// Send the modified response onwards
		return status, headers, body
	}
}
