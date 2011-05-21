package mango

import (
	"http"
)

func prepareSession(cookies []*http.Cookie, key, secret string) map[string]interface{} {
	return make(map[string]interface{})
}

func commitSession(headers Headers, attributes map[string]interface{}, key, secret string) {
}

func Sessions(secret, key string) Middleware {
	return func(env Env, app App) (Status, Headers, Body) {
		env["mango.session"] = prepareSession(env.Request().Cookie, key, secret)
		status, headers, body := app(env)
		commitSession(headers, env["mango.session"].(map[string]interface{}), key, secret)
		return status, headers, body
	}
}
