package mango

import (
	"log"
)

func Logger(logger *log.Logger) Middleware {
	return func(env Env, app App) (Status, Headers, Body) {
		env["mango.logger"] = logger
		return app(env)
	}
}
