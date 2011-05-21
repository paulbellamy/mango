package mango

import (
  "bytes"
  "fmt"
  "http"
  "log"
  "os"
  "strings"
  "template"
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
  Address    string
  middleware []Middleware
  app        App
}

func (this *Mango) Version() []int {
  return []int{0, 1}
}

func (this *Mango) Middleware(middleware ...Middleware) {
  this.middleware = middleware
}

func (this *Mango) buildStack() http.HandlerFunc {
  stack := this.middleware
  compiled_app := bundle(append(stack, middlewareify(this.app))...)
  return func(w http.ResponseWriter, r *http.Request) {
    env := make(map[string]interface{})
    env["REQUEST_METHOD"] = r.Method
    env["REQUEST_PATH"] = r.URL.Path
    env["PATH_INFO"] = r.URL.Path
    env["QUERY_STRING"] = r.URL.RawQuery // failing
    env["SERVER_HOST"] = r.Host
    split_host := strings.Split(r.Host, ":", 2)
    env["SERVER_NAME"] = split_host[0]
    env["SERVER_PORT"] = split_host[1]
    env["REMOTE_ADDR"] = r.RemoteAddr

    if strings.ToUpper(r.URL.Scheme) == "HTTPS" {
      env["HTTPS"] = true
      env["HTTPS_HOST"] = r.Host
      env["HTTPS_USER_AGENT"] = r.UserAgent
      env["HTTPS_COOKIE"] = r.Cookie
    } else {
      env["HTTP"] = true
      env["HTTP_HOST"] = r.Host
      env["HTTP_USER_AGENT"] = r.UserAgent
      env["HTTP_COOKIE"] = r.Cookie
    }

    env["mango.request"] = r
    env["mango.version"] = this.Version()
    env["mango.url_scheme"] = r.URL.Scheme // "http" or "https" depending on URL scheme // failing
    env["mango.input"] = r.Body            // input stream
    //env["mango.errors"] // error stream

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
  if this.Address == "" {
    this.Address = "0.0.0.0:8000"
  }
  log.Println("Starting Mango Server On:", this.Address)
  http.HandleFunc("/", this.buildStack())
  return http.ListenAndServe(this.Address, nil)
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
