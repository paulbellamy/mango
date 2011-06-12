package mango

import (
	"http"
	"http/httptest"
	"testing"
	"runtime"
)

func showErrorsTestServer(env Env) (Status, Headers, Body) {
	panic("foo!")
	return 200, Headers{}, Body("Hello World!")
}

func init() {
	runtime.GOMAXPROCS(4)
}

func TestShowErrors(t *testing.T) {
	// Compile the stack
	showErrorsStack := new(Stack)
	showErrorsStack.Middleware(ShowErrors("<html><body>{Error|html}</body></html>"))
	showErrorsApp := showErrorsStack.Compile(showErrorsTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	status, _, body := showErrorsApp(Env{"mango.request": &Request{request}})

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

	stack := new(Stack)
	stack.Middleware(ShowErrors("<html><body>{Error|html}</body></html>"))
	testServer := httptest.NewServer(stack.HandlerFunc(showErrorsTestServer))
	defer testServer.Close()
	address := testServer.URL

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		http.Get(address)
	}
	b.StopTimer()
}
