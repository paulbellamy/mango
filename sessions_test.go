package mango

import (
	"net/http"
	"strings"
	"testing"
)

func TestSessionEncodingDecoding(t *testing.T) {
	cookie := map[string]interface{}{"value": "foo"}
	secret := "secret"
	result := decodeCookie(encodeCookie(cookie, secret), secret)

	if len(result) != len(cookie) {
		t.Error("Expected cookie to equal:", cookie, "got:", result)
	}

	if result["value"] != cookie["value"] {
		t.Error("Expected cookie[\"value\"] to equal:", cookie["value"], "got:", result["value"])
	}
}

func TestSessions(t *testing.T) {
	sessionsTestServer := func(env Env) (Status, Headers, Body) {
		counter := env.Session()["counter"].(int)
		if counter != 1 {
			t.Error("Expected session[\"counter\"] to equal:", 1, "got:", counter)
		}
		env.Session()["counter"] = counter + 1
		return 200, Headers{}, Body("Hello World!")
	}

	// Compile the stack
	sessionsStack := new(Stack)
	sessionsStack.Middleware(Sessions("my_secret", "my_key", &CookieOptions{Domain: ".my.domain.com"}))
	sessionsApp := sessionsStack.Compile(sessionsTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	initial_cookie := new(http.Cookie)
	initial_cookie.Name = "my_key"
	initial_cookie.Value = encodeCookie(map[string]interface{}{"counter": 1}, "my_secret")
	initial_cookie.Domain = ".my.domain.com"
	request.AddCookie(initial_cookie)
	status, headers, _ := sessionsApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	// base 64 encoded, hmac-hashed, and gob encoded stuff
	cookie := headers.Get("Set-Cookie")
	if cookie == "" {
		t.Error("Expected the Set-Cookie header to be set")
	}

	unparsed := strings.Split(strings.Split(cookie, "=")[1], ";")[0]
	value := decodeCookie(unparsed, "my_secret")
	expected_value := map[string]interface{}{"counter": int(2)}
	if len(value) != len(expected_value) {
		t.Error("Expected cookie to equal:", expected_value, "got:", value)
	}

	if value["counter"].(int) != expected_value["counter"].(int) {
		t.Error("Expected cookie[\"counter\"] to equal:", expected_value["counter"], "got:", value["counter"])
	}

	if !strings.Contains(headers.Get("Set-Cookie"), "; Domain=.my.domain.com") {
		t.Error("Expected cookie ", headers.Get("Set-Cookie"), " to contain: '; Domain=.my.domain.com'")
	}
}

func BenchmarkSessions(b *testing.B) {
	b.StopTimer()

	sessionsTestServer := func(env Env) (Status, Headers, Body) {
		env.Session()["counter"] = env.Session()["counter"].(int) + 1
		return 200, Headers{}, Body("Hello World!")
	}

	sessionsStack := new(Stack)
	sessionsStack.Middleware(Sessions("my_secret", "my_key", &CookieOptions{Domain: ".my.domain.com"}))
	sessionsApp := sessionsStack.Compile(sessionsTestServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sessionsApp(Env{"mango.request": &Request{request}})
	}
	b.StopTimer()
}
