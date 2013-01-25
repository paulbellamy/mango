package mango

import (
	"encoding/base64"
	"errors"
	"strings"
)

func defaultFailure() (Status, Headers, Body) {
	return 401, Headers{"WWW-Authenticate": []string{"Basic realm=\"Basic\""}, "Content-Type": []string{"text/html"}}, Body("Access Denied.") // default failure page
}

func BasicAuth(auth func(string, string, Request, error) bool, failure func(Env) (Status, Headers, Body)) Middleware {
	return func(env Env, app App) (Status, Headers, Body) {

		if auth == nil { // fail auth by default if you use this middleware
			return defaultFailure()
		}

		username, password, err := getAuth(env.Request())

		if auth(username, password, *env.Request(), err) { // check users auth function 
			return app(env)
		}

		if failure == nil { // if no special failure function
			return defaultFailure()
		}

		return failure(env)
	}
}

// get username and password from header
func getAuth(req *Request) (string, string, error) {

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
