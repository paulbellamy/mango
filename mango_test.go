package mango

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func helloWorld(env Env) (Status, Headers, Body) {
	return 200, Headers{}, Body("Hello World!")
}

func init() {
	fmt.Println("Testing Mango Version:", VersionString())
}

func TestHelloWorld(t *testing.T) {
	stack := new(Stack)
	testServer := httptest.NewServer(stack.HandlerFunc(helloWorld))
	defer testServer.Close()

	var client = http.Client{}

	// Request against it
	response, err := client.Get(testServer.URL)

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
	b.StopTimer()

	stack := new(Stack)
	testServer := httptest.NewServer(stack.HandlerFunc(helloWorld))
	defer testServer.Close()
	address := testServer.URL

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		http.Get(address)
	}
	b.StopTimer()
}
