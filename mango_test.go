package mango

import (
	"http"
	"io/ioutil"
	"testing"
	"bytes"
	"fmt"
	"log"
	"runtime"
)

var testServer = Stack{}
var loggerBuffer = &bytes.Buffer{}
var client = http.Client{}

func helloWorld(env Env) (Status, Headers, Body) {
	return 200, Headers{}, Body("Hello World!")
}

func loggerTestServer(env Env) (Status, Headers, Body) {
	env.Logger().Println("Never gonna give you up")
	return 200, Headers{}, Body("Hello World!")
}

func showErrorsTestServer(env Env) (Status, Headers, Body) {
	panic("foo!")
	return 200, Headers{}, Body("Hello World!")
}

func init() {
	runtime.GOMAXPROCS(4)

	fmt.Println("Testing Mango Version:", VersionString())

	testRoutes := make(map[string]App)

	testRoutes["/hello"] = new(Stack).Compile(helloWorld)

	loggerStack := new(Stack)
	custom_logger := log.New(loggerBuffer, "prefixed:", 0)
	loggerStack.Middleware(Logger(custom_logger))
	testRoutes["/logger"] = loggerStack.Compile(loggerTestServer)

	showErrorsStack := new(Stack)
	showErrorsStack.Middleware(ShowErrors("<html><body>{Error|html}</body></html>"))
	testRoutes["/show_errors"] = showErrorsStack.Compile(showErrorsTestServer)

	fullStack := new(Stack)
	fullStackTestRoutes := make(map[string]App)
	fullStackTestRoutes["/full_stack/[0-9]+"] = routingATestServer
	fullStackTestRoutes["/full_stack/[a-z]+"] = routingBTestServer
	fullStack.Middleware(ShowErrors("<html><body>{Error|html}</body></html>"),
		Logger(custom_logger),
		Routing(fullStackTestRoutes))
	testRoutes["/full_stack(.*)"] = fullStack.Compile(helloWorld)

	testServer.Middleware(Routing(testRoutes))
	testServer.Address = "localhost:3000"
	go testServer.Run(helloWorld)
}

func TestHelloWorld(t *testing.T) {
	// Request against it
	response, err := client.Get("http://localhost:3000/hello")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Expected status to equal 200, got:", response.StatusCode)
	}

	body, _ := ioutil.ReadAll(response.Body)
	if string(body) != "Hello World!" {
		t.Error("Expected body:", string(body), "to equal: \"Hello World!\"")
	}
}

func BenchmarkHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/hello")
	}
}

func TestLogger(t *testing.T) {
	// Request against it
	response, err := client.Get("http://localhost:3000/logger")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Expected status to equal 200, got:", response.StatusCode)
	}

	expected := "prefixed:Never gonna give you up\n"
	if loggerBuffer.String() != expected {
		t.Error("Expected logger to print: \"", expected, "\" got: \"", loggerBuffer.String(), "\"")
	}
}

func BenchmarkLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/logger")
	}
}

func TestShowErrors(t *testing.T) {
	// Request against it
	response, err := client.Get("http://localhost:3000/show_errors")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 500 {
		t.Error("Expected status to equal 500, got:", response.StatusCode)
	}

	expected := "<html><body>foo!</body></html>"
	got, _ := ioutil.ReadAll(response.Body)
	if string(got) != expected {
		t.Error("Expected response body to equal: \"", expected, "\" got: \"", string(got), "\"")
	}
}

func BenchmarkShowErrors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/show_errors")
	}
}

func BenchmarkFullStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/full_stack/123")
	}
}
