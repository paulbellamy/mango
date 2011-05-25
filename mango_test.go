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

func init() {
	runtime.GOMAXPROCS(4)

	testRoutes := make(map[string]App)

	testRoutes["/hello"] = new(Stack).Compile(helloWorld)

	loggerStack := new(Stack)
	custom_logger := log.New(loggerBuffer, "prefixed:", 0)
	loggerStack.Middleware(Logger(custom_logger))
	testRoutes["/logger"] = loggerStack.Compile(loggerTestServer)

	testServer.Middleware(Routing(testRoutes))
	testServer.Address = "localhost:3000"
	go testServer.Run(helloWorld)
}

func helloWorld(env Env) (Status, Headers, Body) {
	return 200, make(map[string]string), Body("Hello World!")
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

func loggerTestServer(env Env) (Status, Headers, Body) {
	env.Logger().Println("Never gonna give you up")
	return 200, make(map[string]string), Body("Hello World!")
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
