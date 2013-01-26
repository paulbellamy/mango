package mango

import (
	"net/http"
	"testing"
)

func successPage(env Env) (Status, Headers, Body) {
	return 200, Headers{"Content-Type": []string{"text/html"}}, Body("auth success")
}

func failurePage(env Env) (Status, Headers, Body) {
	return 403, Headers{"Content-Type": []string{"text/html"}}, Body("auth failed")
}

// Example auth function
func auth(username string, password string, req Request, err error) bool {

	if username == "foo" && password == "foo" {
		return true
	}

	return false
}

func TestSuccessAuthRequest(t *testing.T) {

	basicAuthStack := new(Stack)
	basicAuthStack.Middleware(BasicAuth(auth, failurePage))

	basicAuthApp := basicAuthStack.Compile(successPage)

	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	request.SetBasicAuth("foo", "foo")

	status, _, _ := basicAuthApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Request did not succeed, expected status 200, got:", status)
	}
}

func TestFailureAuthRequest(t *testing.T) {

	basicAuthStack := new(Stack)
	basicAuthStack.Middleware(BasicAuth(auth, failurePage))

	basicAuthApp := basicAuthStack.Compile(successPage)

	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	request.SetBasicAuth("fail", "fail")

	status, _, _ := basicAuthApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 403 {
		t.Error("Request did not succeed, expected status 403, got:", status)
	}
}

func TestFailByDefault(t *testing.T) {

	basicAuthStack := new(Stack)
	basicAuthStack.Middleware(BasicAuth(nil, nil))

	basicAuthApp := basicAuthStack.Compile(successPage)

	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)

	status, _, _ := basicAuthApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	// TODO test header
	if status != 401 {
		t.Error("Request did not succeed, expected status 403, got:", status)
	}
}
