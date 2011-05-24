package mango

import (
	"fmt"
	"http"
	"log"
	"os"
)

type Request struct {
	*http.Request
}

type Status int
type Headers map[string]string
type Body string

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

func (this *Stack) Version() []int {
	return []int{0, 1}
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
		env["mango.version"] = this.Version()

		status, headers, body := compiled_app(env)

		for key, value := range headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(int(status))
		fmt.Fprintf(w, string(body))
	}
}

func (this *Stack) Run(app App) os.Error {
	if this.Address == "" {
		this.Address = "0.0.0.0:8000"
	}
	fmt.Println("Starting Mango Stack On:", this.Address)
	http.HandleFunc("/", this.HandlerFunc(app))
	return http.ListenAndServe(this.Address, nil)
}
