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

func sessionsTestServer(env Env) (Status, Headers, Body) {
	env.Session()["test_attribute"] = "Never gonna give you up"
	return 200, Headers{}, Body("Hello World!")
}

func showErrorsTestServer(env Env) (Status, Headers, Body) {
	panic("foo!")
	return 200, Headers{}, Body("Hello World!")
}

func routingATestServer(env Env) (Status, Headers, Body) {
	return 200, Headers{}, Body("Server A")
}

func routingBTestServer(env Env) (Status, Headers, Body) {
	return 200, Headers{}, Body("Server B")
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

	sessionsStack := new(Stack)
	sessionsStack.Middleware(Sessions("my_secret", "my_key", ".my.domain.com"))
	testRoutes["/sessions"] = sessionsStack.Compile(sessionsTestServer)

	showErrorsStack := new(Stack)
	showErrorsStack.Middleware(ShowErrors("<html><body>{Error|html}</body></html>"))
	testRoutes["/show_errors"] = showErrorsStack.Compile(showErrorsTestServer)

	routingStack := new(Stack)
	routingTestRoutes := make(map[string]App)
	routingTestRoutes["/routing/[0-9]+"] = routingATestServer
	routingTestRoutes["/routing/[a-z]+"] = routingBTestServer
	routingStack.Middleware(Routing(routingTestRoutes))
	testRoutes["/routing(.*)"] = routingStack.Compile(helloWorld)

	fullStack := new(Stack)
	fullStackTestRoutes := make(map[string]App)
	fullStackTestRoutes["/full_stack/[0-9]+"] = routingATestServer
	fullStackTestRoutes["/full_stack/[a-z]+"] = routingBTestServer
	fullStack.Middleware(ShowErrors("<html><body>{Error|html}</body></html>"),
		Logger(custom_logger),
		Routing(fullStackTestRoutes))
	testRoutes["/full_stack(.*)"] = fullStack.Compile(helloWorld)

	testServer.Middleware(Static("./static"), Routing(testRoutes))
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

func TestSessions(t *testing.T) {
	// Request against it
	response, err := client.Get("http://localhost:3000/sessions")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Expected status to equal 200, got:", response.StatusCode)
	}

	expected_name := "my_key"
	if response.SetCookie[0].Name != expected_name {
		t.Error("Expected Set-Cookie name to equal: \"", expected_name, "\" got: \"", response.SetCookie[0].Name, "\"")
	}

	// base 64 encoded, hmac-hashed, and gob encoded stuff
	expected_value := "Dv+BBAEC/4IAAQwBEAAANf+CAAEOdGVzdF9hdHRyaWJ1dGUGc3RyaW5nDBkAF05ldmVyIGdvbm5hIGdpdmUgeW91IHVw--bdHyJ5lvPpk6EoZiSSSiHKZtQHk="
	if response.SetCookie[0].Value != expected_value {
		t.Error("Expected Set-Cookie value to equal: \"", expected_value, "\" got: \"", response.SetCookie[0].Value, "\"")
	}

	expected_domain := ".my.domain.com"
	if response.SetCookie[0].Domain != expected_domain {
		t.Error("Expected Set-Cookie domain to equal: \"", expected_domain, "\" got: \"", response.SetCookie[0].Domain, "\"")
	}
}

func BenchmarkSessions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/sessions")
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

func TestRouting(t *testing.T) {
	// Request server a
	response, err := client.Get("http://localhost:3000/routing/123")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Expected status to equal 200, got:", response.StatusCode)
	}

	expected := "Server A"
	body, _ := ioutil.ReadAll(response.Body)
	if string(body) != expected {
		t.Error("Expected response body to equal: \"", expected, "\" got: \"", string(body), "\"")
	}

	// Request server b
	response, err = client.Get("http://localhost:3000/routing/abc")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Expected status to equal 200, got:", response.StatusCode)
	}

	expected = "Server B"
	body, _ = ioutil.ReadAll(response.Body)
	if string(body) != expected {
		t.Error("Expected response body to equal: \"", expected, "\" got: \"", string(body), "\"")
	}
}

func BenchmarkRouting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/routing/123")
	}
}

func TestStatic(t *testing.T) {
	// Request against it
	response, err := client.Get("http://localhost:3000/static.html")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Expected status to equal 200, got:", response.StatusCode)
	}

	body, _ := ioutil.ReadAll(response.Body)
	expected := "<h1>I'm a static test file</h1>\n"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal: \"", expected, "\"")
	}
}

func BenchmarkStatic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/static.html")
	}
}

func BenchmarkFullStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/full_stack/123")
	}
}
