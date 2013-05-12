package mango

import (
	"net/http"
	"testing"
)

func routingTestServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func routingATestServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server A"))
}

func routingBTestServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server B"))
}

func TestRoutesSuccess(t *testing.T) {
	// Compile the stack
	app := Routes(
		GET("/a", routingATestServer),
		GET("/b", routingBTestServer),
		ANY(".*", routingTestServer),
	)

	// Request against A
	request, err := http.NewRequest("GET", "http://localhost:3000/a", nil)
	response := NewBufferedResponseWriter(nil)
	app(response, request)
	status := response.Status
	body := response.Body.String()

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "Server A"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}

	// Request against B
	request, err = http.NewRequest("GET", "http://localhost:3000/b", nil)
	response = NewBufferedResponseWriter(nil)
	app(response, request)
	status = response.Status
	body = response.Body.String()

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected = "Server B"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestRoutesFailure(t *testing.T) {
	// Compile the stack
	app := Routes(
		GET("/a", routingATestServer),
		GET("/b", routingBTestServer),
		ANY(".*", routingTestServer),
	)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/c", nil)
	response := NewBufferedResponseWriter(nil)
	app(response, request)
	status := response.Status
	body := response.Body.String()

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "Hello World!"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func BenchmarkRoutes(b *testing.B) {
	b.StopTimer()

	app := Routes(
		GET("/a", routingATestServer),
		GET("/b", routingBTestServer),
		ANY(".*", routingTestServer),
	)

	request, _ := http.NewRequest("GET", "http://localhost:3000/a", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		response := NewBufferedResponseWriter(nil)
		app(response, request)
	}
	b.StopTimer()
}
