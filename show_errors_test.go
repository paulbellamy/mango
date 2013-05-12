package mango

import (
	"net/http"
	"testing"
)

func showErrorsTestServer(w http.ResponseWriter, r *http.Request) {
	panic("foo!")
	w.Write([]byte("Hello World!"))
}

func TestShowErrors(t *testing.T) {
	// Compile the app
	app := ShowErrors("<html><body>{{.Error|html}}</body></html>", showErrorsTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	response := &MockResponseWriter{}
	app(response, request)
	status := response.Status
	body := response.Body.String()

	if err != nil {
		t.Error(err)
	}

	if status != 500 {
		t.Error("Expected status to equal 500, got:", status)
	}

	expected := "<html><body>foo!</body></html>"
	if string(body) != expected {
		t.Error("Expected response body to equal: \"", expected, "\" got: \"", string(body), "\"")
	}
}

func BenchmarkShowErrors(b *testing.B) {
	b.StopTimer()

	app := ShowErrors("<html><body>{Error|html}</body></html>", showErrorsTestServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/", nil)
	response := &MockResponseWriter{}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		app(response, request)
		response.Status = 0
		response.Body.Reset()
	}
	b.StopTimer()
}
