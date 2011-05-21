package main

import (
	"bytes"
	"fmt"
	"http"
	"log"
	"os"
	"template"
	"time"
)

type Env map[string]interface{}
type Status int
type Headers map[string]string
type Body string

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

type Mango struct {
	address    string
	middleware []Middleware
	app        App
}

func (this *Mango) Middleware(middleware ...Middleware) {
	this.middleware = middleware
}

func (this *Mango) buildStack() http.HandlerFunc {
	stack := this.middleware
	compiled_app := bundle(append(stack, middlewareify(this.app))...)
	return func(w http.ResponseWriter, r *http.Request) {
		env := make(map[string]interface{})
		env["mango.request"] = r
		status, headers, body := compiled_app(env)
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(int(status))
		fmt.Fprintf(w, string(body))
	}
}

func (this *Mango) Run(app App) os.Error {
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

// An example of how to pass runtime config to Middleware
func Logger(prefix string) Middleware {
	return func(env Env, app App) (Status, Headers, Body) {
		status, headers, body := app(env)
		log.Println(prefix, env["mango.request"].(*http.Request).Method, env["mango.request"].(*http.Request).RawURL, status)
		return status, headers, body
	}
}

func ShowErrors(templateString string) Middleware {
	if templateString == "" {
		templateString = `
      <html>
      <body>
        <p>
          {Error|html}
        </p>
      </body>
      </html>
    `
	}

	errorTemplate := template.MustParse(templateString, nil)

	return func(env Env, app App) (status Status, headers Headers, body Body) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Error: ", err)
				buffer := bytes.NewBufferString("")
				errorTemplate.Execute(buffer, struct{ Error string }{err.(string)})
				status = 500
				headers = make(map[string]string)
				body = Body(buffer.String())
			}
		}()

		return app(env)
	}
}

func Hello(env Env) (Status, Headers, Body) {
	return 200, map[string]string{"Never-Gonna": "Give you up!"}, Body(fmt.Sprintf("%d", time.Seconds()))
}

func main() {
	mango := new(Mango)
	mango.address = ":3000"
	mango.Middleware(Logger("my_custom_prefix:"), ShowErrors(""))
	mango.Run(Hello)
}
