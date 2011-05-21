package mango

import ()

func prepareSession(env Env, key, secret string) {
}

func commitSession(headers Headers, env Env, key, secret string) {
}

func Sessions(secret, key string) Middleware {
	return func(env Env, app App) (Status, Headers, Body) {
		prepareSession(env, key, secret)
		status, headers, body := app(env)
		commitSession(headers, env, key, secret)
		return status, headers, body
	}
}
