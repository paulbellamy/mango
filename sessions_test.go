package mango

import (
	"http"
	"runtime"
	"testing"
)

func sessionsTestServer(env Env) (Status, Headers, Body) {
	env.Session()["test_attribute"] = "Never gonna give you up"
	return 200, Headers{}, Body("Hello World!")
}

func init() {
	runtime.GOMAXPROCS(4)
}

func TestSessions(t *testing.T) {
	// Compile the stack
	sessionsStack := new(Stack)
	sessionsStack.Middleware(Sessions("my_secret", "my_key", ".my.domain.com"))
	sessionsApp := sessionsStack.Compile(sessionsTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	status, headers, _ := sessionsApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	// base 64 encoded, hmac-hashed, and gob encoded stuff
	cookie := headers.Get("Set-Cookie")
	expected_cookie := "my_key=Dv+BBAEC/4IAAQwBEAAANf+CAAEOdGVzdF9hdHRyaWJ1dGUGc3RyaW5nDBkAF05ldmVyIGdvbm5hIGdpdmUgeW91IHVw--bdHyJ5lvPpk6EoZiSSSiHKZtQHk=; Domain=.my.domain.com;"
	if cookie != expected_cookie {
		t.Error("Expected Set-Cookie to equal: \"", expected_cookie, "\" got: \"", cookie, "\"")
	}
}

func BenchmarkSessions(b *testing.B) {
	b.StopTimer()

	sessionsStack := new(Stack)
	sessionsStack.Middleware(Sessions("my_secret", "my_key", ".my.domain.com"))
	sessionsApp := sessionsStack.Compile(sessionsTestServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sessionsApp(Env{"mango.request": &Request{request}})
	}
	b.StopTimer()
}
