# Mango

Mango is a modular web-application library for Go. It has wrappers for http.HandlerFuncs, and helper functions for common web operations. Mango aims to work seamlessly with the existing Go http library as much as possible.


## API

As in Go's http package, the API is very minimal.

Applications are of the form:

    func Hello(w http.ResponseWriter, r *http.Response) {
      w.Write([]byte("Hello World))
    }

Wrappers are designed to do pre/post processing around an application. For example, handling JSONP requests.

There are also several helpers designed to simplify common tasks in a web application. For example, session handling, and basic auth.

Where possible everything is designed to work independently, but also be composable.


## Installation

   $ go get github.com/paulbellamy/mango


## Available Wrappers
  
* ShowErrors

  Usage: mango.ShowErrors(templateString string, app http.HandlerFunc) http.HandlerFunc

  Catch any panics thrown from the app, and display them in an HTML template. If templateString is "", a default template is used. Not recommended to use the default template in production as it could provide information helpful to attackers.

* Routes

  Usage: mango.Routes(handlers ...http.HandlerFunc) http.HandlerFunc

  Takes a list of http.HandlerFuncs, and tries each one successively until one writes some data to the ResponseWriter.

* GET/POST/PUT/DELETE/HEAD/OPTIONS/PATCH/ANY

  Usage: mango.GET(path string, http.HandlerFunc) http.HandlerFunc

  Creates a scoped handler which will only respond to the matching path regexp, and HTTP method. Best used in conjunction with the Routes helper.

* Static

  Usage: mango.Static(root http.FileSystem, app http.HandlerFunc) http.HandlerFunc

  Serves static files from the directory provided. If file is not found, it calls app. As opposed to the built-in FileServer in net/http, this one will not return a 404 if an app is specified.

* JSONP

  Usage: mango.JSONP(app http.HandlerFunc) http.HandlerFunc

  Provides JSONP support. If a request has a 'callback' parameter, and your application responds with a Content-Type of "application/json", the JSONP middleware will wrap the response in the callback function and set the Content-Type to "application/javascript".

## Available Helpers

* NewBufferedResponseWriter(w http.ResponseWriter) *BufferedResponseWriter

  Creates a new BufferedResponseWriter, which wraps a regular http.ResponseWriter. Data can be written to the BufferedResponseWriter until finished, then calling Flush() will send it to the client. If no http.ResponseWriter is wrapped, the BufferedResponseWriter will still function normally aside from the Flush() method.

* Session

  Usage: mango.Session(*r http.Request)

  Basic session management. Returns a SessionWrapper object representing the session. Data stored here can be stored into the session cookie by calling the Write(w http.ResponseWriter) method.

* Basic Auth

  Usage: mango.BasicAuth(r *http.Request) (username string, password string)

  Fetches the basic auth username and password from an http.Request object.

## Example App

    package main

    import (
      . "github.com/paulbellamy/mango"
      "log"
      "net/http"
    )

    func Hello(w http.ResponseWriter, r *http.Request) {
      log.Println("Got a", r.Method, "request for", r.RequestURI)
      return w.Write([]byte("Hello World!"))
    }

    func main() {
      app := ShowErrors("", Hello)
      http.HandleFunc("/", app)
      http.ListenAndServe(":3000", nil)
    }

## Routing example

The following example routes "/hello" traffic to the hello handler and
"/bye" traffic to the bye handler, any other traffic goes to
routeNotFound handler returning a 404.

    package main

    import (
      . "github.com/paulbellamy/mango"
      "log"
      "net/http"
    )

    func hello(w http.ResponseWriter, r *http.Request) {
      log.Println("Got a", env.Request().Method, "request for", env.Request().RequestURI)
      w.Write([]byte("Hello World!"))
    }

    func bye(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("Bye Bye!"))
    }

    func main() {
      app :=
        ShowErrors("<html><body>{Error|html}</body></html>",
          Routes(
            GET("/hello", hello),
            GET("/bye", bye),
            DEFAULT(http.NotFound),
          ))

      http.HandleFunc("/", app)
      http.ListenAndServe(":3000", nil)
      log.Println("Running server on: :3000")
    }


## About

Mango was written by [Paul Bellamy](http://paulbellamy.com). 

Follow me on [Twitter](http://www.twitter.com/pyrhho)!
