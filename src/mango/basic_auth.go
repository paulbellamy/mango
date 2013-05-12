package mango

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

// get username and password from header
func BasicAuth(req *http.Request) (string, string, error) {

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
