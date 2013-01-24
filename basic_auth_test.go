package mango

import (
	"net/http"
	"testing"
)

func successPage(env Env) (Status, Headers, Body) {
	return 200, Headers{"Content-Type": []string{"text/html"}}, Body("auth success")
}

func auth(username string, password string, req Request, err error) bool {

	if username == "foo" && password == "foo" { 
		return true
	}

	return false
}

var test *testing.T
func TestAuthRequest(t *testing.T) {

	test = t
	basicAuthStack := new(Stack)
	basicAuthStack.Middleware(BasicAuth(auth))

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



//func TestBasicRequest(t *testing.T) {
//
//	basicAuthStack := new(Stack)
//	basicAuthStack.Middleware(BasicAuth)
//
//	basicAuthApp := basicAuthStack.Compile(successPage)
//	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
//	status, _, _ := basicAuthApp(Env{"mango.request": &Request{request}})
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	if status != 401 {
//		t.Error("Expected status to equal 401, got:", status)
//	}
//}
//
