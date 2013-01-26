// Mango is a modular web-application framework for Go, inspired by Rack and PEP333.
package mango

import (
	"fmt"
	"log"
	"net/http"
	"net/textproto"
	"os"
)

type Request struct {
	*http.Request
}

type Status int
type Body string

type Headers http.Header

func (h Headers) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}

func (h Headers) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

func (h Headers) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

func (h Headers) Del(key string) {
	textproto.MIMEHeader(h).Del(key)
}

type Env map[string]interface{}

func (this Env) Logger() *log.Logger {
	if this["mango.logger"] == nil {
		this["mango.logger"] = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}

	return this["mango.logger"].(*log.Logger)
}

func (this Env) Request() *Request {
	return this["mango.request"].(*Request)
}

func (this Env) Session() map[string]interface{} {
	return this["mango.session"].(map[string]interface{})
}

// This is the core app the user has written
type App func(Env) (Status, Headers, Body)

// These are pieces of middleware,
// which 'wrap' around the core App
// (and each other)
type Middleware func(Env, App) (Status, Headers, Body)

// Bundle a given list of Middleware pieces into a App
func bundle(r ...Middleware) App {
	if len(r) <= 1 {
		// Terminate the innermost piece of Middleware
		// Basically stops it from recursing any further.
		return func(input Env) (Status, Headers, Body) {
			return r[0](input, func(Env) (Status, Headers, Body) {
				panic("Core Mango App should never call it's upstream function.")
			})
		}
	}
	return wrap(r[0], bundle(r[1:]...))
}

// Attach a piece of Middleware to the outside
// of a App. This wraps the inner App
// inside the outer Middleware.
func wrap(middleware Middleware, app App) App {
	return func(input Env) (Status, Headers, Body) {
		return middleware(input, app)
	}
}

// Convert a App into Middleware
// We convert the core app into a Middleware
// so we can pass it to Bundle as part of the
// stack. Because the App does not call its
// upstream method, the resulting Middleware
// will just ignore any upstream passed to it.
func middlewareify(app App) Middleware {
	return func(input Env, upstream App) (Status, Headers, Body) {
		return app(input)
	}
}

type Stack struct {
	Address    string
	middleware []Middleware
	app        App
}

func Version() []int {
	return []int{0, 5, 0}
}

func VersionString() string {
	v := Version()
	return fmt.Sprintf("%d.%02d.%02d", v[0], v[1], v[2])
}

func (this *Stack) Middleware(middleware ...Middleware) {
	this.middleware = middleware
}

func (this *Stack) Compile(app App) App {
	this.app = app
	return bundle(append(this.middleware, middlewareify(this.app))...)
}

func (this *Stack) HandlerFunc(app App) http.HandlerFunc {
	compiled_app := this.Compile(app)
	return func(w http.ResponseWriter, r *http.Request) {
		env := make(map[string]interface{})
		env["mango.request"] = &Request{r}
		env["mango.version"] = Version()

		status, headers, body := compiled_app(env)

		for key, values := range headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(int(status))
		w.Write([]byte(body))
	}
}

func (this *Stack) Run(app App) error {
	if this.Address == "" {
		this.Address = "0.0.0.0:8000"
	}
	fmt.Println("Starting Mango Stack On:", this.Address)
	http.HandleFunc("/", this.HandlerFunc(app))
	return http.ListenAndServe(this.Address, nil)
}
