package mango

import (
	"http"
	"testing"
	"fmt"
	"runtime"
)

func jsonServer(env Env) (Status, Headers, Body) {
	return 200, Headers{"Content-Type": []string{"application/json"}}, Body("{\"foo\":\"bar\"}")
}

func nonJsonServer(env Env) (Status, Headers, Body) {
	return 200, Headers{"Content-Type": []string{"text/html"}}, Body("<h1>Hello World!</h1>")
}

func init() {
	runtime.GOMAXPROCS(4)

	fmt.Println("Testing Mango-JSONP Version:", VersionString())
}

func TestJSONPSuccess(t *testing.T) {
	// Compile the stack
	jsonpStack := new(Stack)
	jsonpStack.Middleware(JSONP)
	jsonpApp := jsonpStack.Compile(jsonServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/?callback=parseResponse", nil)
	status, headers, body := jsonpApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	if headers.Get("Content-Type") != "application/javascript" {
		t.Error("Expected Content-Type to equal \"application/javascript\", got:", headers.Get("Content-Type"))
	}

	expected := "parseResponse({\"foo\":\"bar\"})"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestNonJSONPSuccess(t *testing.T) {
	// Compile the stack
	nonJsonpStack := new(Stack)
	nonJsonpStack.Middleware(JSONP)
	nonJsonpApp := nonJsonpStack.Compile(nonJsonServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/?callback=parseResponse", nil)
	status, headers, body := nonJsonpApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	if headers.Get("Content-Type") != "text/html" {
		t.Error("Expected Content-Type to equal \"text/html\", got:", headers.Get("Content-Type"))
	}

	expected := "<h1>Hello World!</h1>"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestJSONPNoCallback(t *testing.T) {
	// Compile the stack
	jsonpStack := new(Stack)
	jsonpStack.Middleware(JSONP)
	jsonpApp := jsonpStack.Compile(jsonServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	status, headers, body := jsonpApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	if headers.Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type to equal \"application/json\", got:", headers.Get("Content-Type"))
	}

	expected := "{\"foo\":\"bar\"}"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestJSONPInvalidCallback(t *testing.T) {
	// Compile the stack
	jsonpStack := new(Stack)
	jsonpStack.Middleware(JSONP)
	jsonpApp := jsonpStack.Compile(jsonServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/?callback=invalid(callback)", nil)
	status, headers, body := jsonpApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 400 {
		t.Error("Expected status to equal 400, got:", status)
	}

	if headers.Get("Content-Type") != "text/plain" {
		t.Error("Expected Content-Type to equal \"text/plain\", got:", headers.Get("Content-Type"))
	}

	expected := "Bad Request"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}
