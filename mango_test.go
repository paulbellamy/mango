package mango

import (
	"http"
	"io/ioutil"
	"testing"
	"./mango"
)

func helloWorld(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
	return 200, make(map[string]string), mango.Body("Hello World!")
}

func TestHelloWorld(t *testing.T) {
	// Start up the server
	stack := new(mango.Stack)
	stack.Address = "localhost:3000"
	go stack.Run(helloWorld)

	// Request against it
	client := new(http.Client)
	response, _, err := client.Get("http://localhost:3000/")
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
