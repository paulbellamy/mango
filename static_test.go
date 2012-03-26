package mango

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func staticTestServer(env Env) (Status, Headers, Body) {
	return 200, Headers{"Content-Type": []string{"text/html"}}, Body("<h1>Hello World!</h1>")
}

func TestStaticSuccess(t *testing.T) {
	// Compile the stack
	staticStack := new(Stack)
	staticStack.Middleware(Static("./static"))
	staticApp := staticStack.Compile(staticTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/static.html", nil)
	status, headers, body := staticApp(Env{"mango.request": &Request{request}})

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

	expected = "text/html"
	got := headers.Get("Content-Type")
	if got != expected {
		t.Error("Expected Content-Type:", got, "to equal:", expected)
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

func TestStaticBinaryFile(t *testing.T) {
	// Compile the stack
	staticStack := new(Stack)
	staticStack.Middleware(Static("./static"))
	staticApp := staticStack.Compile(staticTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/binary_file.png", nil)
	status, _, body := staticApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected, err := ioutil.ReadFile("./static/binary_file.png")
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare([]byte(body), []byte(expected)) != 0 {
		t.Error("Expected body to equal ./static/binary_file.png")
	}
}

func BenchmarkStatic(b *testing.B) {
	b.StopTimer()

	staticStack := new(Stack)
	staticStack.Middleware(Static("./static"))
	staticApp := staticStack.Compile(staticTestServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/static.html", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		staticApp(Env{"mango.request": &Request{request}})
	}
	b.StopTimer()
}
