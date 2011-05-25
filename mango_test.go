package mango

import (
	"http"
	"io/ioutil"
	"testing"
	"bytes"
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

func init() {
	runtime.GOMAXPROCS(4)

	testRoutes := make(map[string]App)

	testRoutes["/hello"] = new(Stack).Compile(helloWorld)

	loggerStack := new(Stack)
	custom_logger := log.New(loggerBuffer, "prefixed:", 0)
	loggerStack.Middleware(Logger(custom_logger))
	testRoutes["/logger"] = loggerStack.Compile(loggerTestServer)

	sessionsStack := new(Stack)
	sessionsStack.Middleware(Sessions("my_secret", "my_key", ".my.domain.com"))
	testRoutes["/sessions"] = sessionsStack.Compile(sessionsTestServer)

	testServer.Middleware(Routing(testRoutes))
	testServer.Address = "localhost:3000"
	go testServer.Run(helloWorld)
}

func TestHelloWorld(t *testing.T) {
	// Request against it
	response, _, err := client.Get("http://localhost:3000/hello")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Expected status to equal 200, got:", response.StatusCode)
	}

	body, _ := ioutil.ReadAll(response.Body)
	if string(body) != "Hello World!" {
		t.Error("Expected body:", body, "to equal: \"Hello World!\"")
	}
}

func BenchmarkHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost:3000/hello")
	}
}

func TestLogger(t *testing.T) {
	// Request against it
	response, _, err := client.Get("http://localhost:3000/logger")

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
	response, _, err := client.Get("http://localhost:3000/sessions")

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

	// base 64 encoded, hashed, and gob encoded stuff
	expected_value := "Dv+BBAEC/4IAAQwBEAAANf+CAAEOdGVzdF9hdHRyaWJ1dGUGc3RyaW5nDBkAF05ldmVyIGdvbm5hIGdpdmUgeW91IHVw--q0x2Xt9XBekiKpL2/MlQ50TcOqg="
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
