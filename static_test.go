package mango

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func staticTestServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Hello World!</h1>"))
}

func TestStaticSuccess(t *testing.T) {
	// Compile the stack
	app := Static("./static", staticTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/static.html", nil)
	response := NewMockResponseWriter()
	app(response, request)
	status := response.Status
	headers := response.Header()
	body := response.Body.String()

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "<h1>I'm a static test file</h1>\n"
	if body != expected {
		t.Error("Expected body:", body, "to equal:", expected)
	}

	expected = "text/html; charset=utf-8"
	got := headers.Get("Content-Type")
	if got != expected {
		t.Error("Expected Content-Type:", got, "to equal:", expected)
	}
}

func TestStaticNoUpstream(t *testing.T) {
	// Compile the stack
	app := Static("./static", nil)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/not_a_file.html", nil)
	response := NewMockResponseWriter()
	app(response, request)
	status := response.Status
	body := response.Body.String()

	if err != nil {
		t.Error(err)
	}

	if status != 0 {
		t.Error("Expected status to equal 0, got:", status)
	}

	expected := ""
	if body != expected {
		t.Error("Expected body:", body, "to equal:", expected)
	}
}

func TestStaticFail(t *testing.T) {
	// Compile the stack
	app := Static("./static", staticTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/not_a_file.html", nil)
	response := NewMockResponseWriter()
	app(response, request)
	status := response.Status
	body := response.Body.String()

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "<h1>Hello World!</h1>"
	if body != expected {
		t.Error("Expected body:", body, "to equal:", expected)
	}
}

func TestStaticBinaryFile(t *testing.T) {
	// Compile the stack
	app := Static("./static", staticTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/binary_file.png", nil)
	response := NewMockResponseWriter()
	app(response, request)
	status := response.Status
	body := response.Body.String()

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

	app := Static("./static", staticTestServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/static.html", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		response := NewMockResponseWriter()
		app(response, request)
	}
	b.StopTimer()
}
