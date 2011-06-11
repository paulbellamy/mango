package mango

import (
	"http"
	"io/ioutil"
	"testing"
	"fmt"
	"runtime"
)

var testServer = Stack{}
var client = http.Client{}

func helloWorld(env Env) (Status, Headers, Body) {
	return 200, Headers{}, Body("Hello World!")
}

func init() {
	runtime.GOMAXPROCS(4)

	fmt.Println("Testing Mango Version:", VersionString())

	testRoutes := make(map[string]App)

	testRoutes["/hello"] = new(Stack).Compile(helloWorld)

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
