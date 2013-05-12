package mango

import (
	"net/http"
	"testing"
)

func jsonServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", "13")
	w.Write([]byte("{\"foo\":\"bar\"}"))
}

func nonJsonServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Hello World!</h1>"))
}

func TestJSONPSuccess(t *testing.T) {
	// Compile the stack
	app := JSONP(jsonServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/?callback=parseResponse", nil)
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

	if headers.Get("Content-Type") != "application/javascript" {
		t.Error("Expected Content-Type to equal \"application/javascript\", got:", headers.Get("Content-Type"))
	}

	if headers.Get("Content-Length") != "28" {
		t.Error("Expected Content-Length to equal \"28\", got:", headers.Get("Content-Length"))
	}

	expected := "parseResponse({\"foo\":\"bar\"})"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestNonJSONPSuccess(t *testing.T) {
	// Compile the stack
	app := JSONP(nonJsonServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/?callback=parseResponse", nil)
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

	if headers.Get("Content-Type") != "text/html" {
		t.Error("Expected Content-Type to equal \"text/html\", got:", headers.Get("Content-Type"))
	}

	if headers.Get("Content-Length") != "" {
		t.Error("Expected Content-Length to equal \"\", got:", headers.Get("Content-Length"))
	}

	expected := "<h1>Hello World!</h1>"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestJSONPNoCallback(t *testing.T) {
	// Compile the stack
	app := JSONP(jsonServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
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

	if headers.Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type to equal \"application/json\", got:", headers.Get("Content-Type"))
	}

	if headers.Get("Content-Length") != "13" {
		t.Error("Expected Content-Length to equal \"13\", got:", headers.Get("Content-Length"))
	}

	expected := "{\"foo\":\"bar\"}"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestJSONPInvalidCallback(t *testing.T) {
	// Compile the stack
	app := JSONP(jsonServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/?callback=invalid(callback)", nil)
	response := NewMockResponseWriter()
	app(response, request)
	status := response.Status
	headers := response.Header()
	body := response.Body.String()

	if err != nil {
		t.Error(err)
	}

	if status != 400 {
		t.Error("Expected status to equal 400, got:", status)
	}

	if headers.Get("Content-Type") != "text/plain" {
		t.Error("Expected Content-Type to equal \"text/plain\", got:", headers.Get("Content-Type"))
	}

	if headers.Get("Content-Length") != "11" {
		t.Error("Expected Content-Length to equal \"11\", got:", headers.Get("Content-Length"))
	}

	expected := "Bad Request"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func BenchmarkJSONP(b *testing.B) {
	b.StopTimer()

	app := JSONP(jsonServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/?callback=parseResponse", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		response := NewMockResponseWriter()
		app(response, request)
	}
	b.StopTimer()
}

func BenchmarkNonJSONP(b *testing.B) {
	b.StopTimer()

	app := JSONP(nonJsonServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/?callback=parseResponse", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		response := NewMockResponseWriter()
		app(response, request)
	}
	b.StopTimer()
}

func BenchmarkJSONPNoCallback(b *testing.B) {
	b.StopTimer()

	app := JSONP(jsonServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		response := NewMockResponseWriter()
		app(response, request)
	}
	b.StopTimer()
}

func BenchmarkJSONPInvalidCallback(b *testing.B) {
	b.StopTimer()

	app := JSONP(jsonServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/?callback=invalid(callback)", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		response := NewMockResponseWriter()
		app(response, request)
	}
	b.StopTimer()
}
