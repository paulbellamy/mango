package main

import (
	"fmt"
	"http"
	"log"
	"os"
	"time"
)

type Env map[string]interface{}
type Status int
type Headers map[string]string
type Body string

// This is the core app the user has written
type MangoApp func(Env) (Status, Headers, Body)

// These are pieces of middleware,
// which 'wrap' around the core MangoApp
// (and each other)
type MangoMiddleware func(Env, MangoApp) (Status, Headers, Body)

// Bundle a given list of MangoMiddleware pieces into a MangoApp
func bundle(r ...MangoMiddleware) MangoApp {
	if len(r) <= 1 {
		// Terminate the innermost piece of MangoMiddleware
		// Basically stops it from recursing any further.
		return func(input Env) (Status, Headers, Body) {
			return r[0](input, func(Env) (Status, Headers, Body) {
				panic("Core Mango App should never call it's upstream function.")
			})
		}
	}
	return wrap(r[0], bundle(r[1:]...))
}

// Attach a piece of MangoMiddleware to the outside
// of a MangoApp. This wraps the inner MangoApp
// inside the outer MangoMiddleware.
func wrap(middleware MangoMiddleware, app MangoApp) MangoApp {
	return func(input Env) (Status, Headers, Body) {
		return middleware(input, app)
	}
}

// Convert a MangoApp into MangoMiddleware
// We convert the core app into a MangoMiddleware
// so we can pass it to Bundle as part of the
// stack. Because the MangoApp does not call its
// upstream method, the resulting MangoMiddleware
// will just ignore any upstream passed to it.
func middlewareify(app MangoApp) MangoMiddleware {
	return func(input Env, upstream MangoApp) (Status, Headers, Body) {
		return app(input)
	}
}

type Mango struct {
	address    string
	middleware []MangoMiddleware
	app        MangoApp
}

func (this *Mango) Middleware(middleware ...MangoMiddleware) {
	this.middleware = middleware
}

func (this *Mango) buildStack() http.HandlerFunc {
	stack := this.middleware
	compiled_app := bundle(append(stack, middlewareify(this.app))...)
	return func(w http.ResponseWriter, r *http.Request) {
		env := make(map[string]interface{})
		status, headers, body := compiled_app(env)
		w.WriteHeader(int(status))
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		fmt.Fprintf(w, string(body))
	}
}

func (this *Mango) Run(app MangoApp) os.Error {
	this.app = app
	if this.address == "" {
		this.address = "0.0.0.0:8000"
	}
	log.Println("Starting Mango Server On:", this.address)
	http.HandleFunc("/", this.buildStack())
	return http.ListenAndServe(this.address, nil)
}


/*************************************
 * End Mango Source
 * Begin Example Usage
 ************************************/

func Logger(env Env, app MangoApp) (Status, Headers, Body) {
	status, headers, body := app(env)
	log.Println(env["REQUEST_METHOD"], env["REQUEST_PATH"], status)
	return status, headers, body
}

func Hello(Env) (Status, Headers, Body) {
	return 200, map[string]string{"Content-Type": "text/html"}, Body(fmt.Sprintf("%d", time.Seconds()))
}

func main() {
	mango := new(Mango)
	mango.address = ":3000"
	mango.Middleware(Logger)
	mango.Run(Hello)
}
