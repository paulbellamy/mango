package mango

import (
	"encoding/base64"
	"strings"
	"errors"
)

func BasicAuth(auth func(string, string, Request, error) bool) Middleware {
	return func(env Env, app App) (Status, Headers, Body) {
		
		username, password, err := deBase64(env.Request())

		if auth(username, password, *env.Request(), err) { // check users auth function 
			return app(env)
		}

		return 401, Headers{"WWW-Authenticate": []string{"Basic Realm=\"Login Required\""}}, Body("need auth")
	}
}

func deBase64(req *Request) (string, string, error) {

	auth64 := req.Header.Get("Authorization")

	if auth64 == "" {
		return "", "", errors.New("No Authorization Header")
	}

	auth, err := base64.StdEncoding.DecodeString(strings.Replace(auth64, "Basic ", "", 1))

	if err != nil {
		return "", "", err 
	}

	result := strings.Split(string(auth), ":")

	return result[0], result[1], nil
}

