package mango

import (
	"net/http"
	"testing"
)

func TestSuccessAuthRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	request.SetBasicAuth("foo", "foo_pass")

  username, password, err := BasicAuth(request)

	if err != nil {
		t.Error(err)
	}

	if username != "foo" || password != "foo_pass" {
    t.Error("Request did not succeed, expected username: foo and password: foo_pass, but got:", username, "and:", password)
	}
}

func TestNoAuthRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
  username, password, err := BasicAuth(request)

	if err == nil {
    t.Error("Expected an error, but got none")
	}

	if username != "" || password != "" {
    t.Error("Request did not succeed, expected username: nil and password: nil, but got:", username, "and:", password)
	}
}
