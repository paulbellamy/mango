package mango

import (
	"http"
	"testing"
	"runtime"
)

func staticTestServer(env Env) (Status, Headers, Body) {
	return 200, Headers{"Content-Type": []string{"text/html"}}, Body("<h1>Hello World!</h1>")
}

func init() {
	runtime.GOMAXPROCS(4)
}

func TestStaticSuccess(t *testing.T) {
	// Compile the stack
	staticStack := new(Stack)
	staticStack.Middleware(Static("./static"))
	staticApp := staticStack.Compile(staticTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/static.html", nil)
	status, _, body := staticApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "<h1>I'm a static test file</h1>\n"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestStaticFail(t *testing.T) {
	// Compile the stack
	staticStack := new(Stack)
	staticStack.Middleware(Static("./static"))
	staticApp := staticStack.Compile(staticTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/not_a_file.html", nil)
	status, _, body := staticApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "<h1>Hello World!</h1>"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}
