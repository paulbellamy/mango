package mango

import (
	"net/http"
	"testing"
)

func redirectTestServer(env Env) (Status, Headers, Body) {
	return Redirect(302, "/somewhere")
}

func TestRedirect(t *testing.T) {
	// Compile the stack
	redirectApp := new(Stack).Compile(redirectTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	status, headers, body := redirectApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 302 {
		t.Error("Expected status to equal 302, got:", status)
	}

	expected := "/somewhere"
	if headers["Location"][0] != expected {
		t.Error("Expected Location header to be: \"", expected, "\" got: \"", headers["Location"][0], "\"")
	}

	expectedBody := Body("")
	if body != expectedBody {
		t.Error("Expected body to be: \"", expected, "\" got: \"", body, "\"")
	}
}

func BenchmarkRedirect(b *testing.B) {
	b.StopTimer()

	redirectApp := new(Stack).Compile(redirectTestServer)

	// Request against it
	request, _ := http.NewRequest("GET", "http://localhost:3000/", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		redirectApp(Env{"mango.request": &Request{request}})
	}
	b.StopTimer()
}
