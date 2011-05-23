package mango

import (
	"regexp"
)

func Routing(routes map[string]App) Middleware {
	return func(env Env, app App) (Status, Headers, Body) {
		path := []byte(env.Request().URL.Path)
		for regex, stack := range routes {
			if matched, err := regexp.Match(regex, path); matched && err == nil {
				// Matched a route return it
				return stack(env)
			}
		}

		// didn't match any of the other routes. pass upstream.
		return app(env)
	}
}
