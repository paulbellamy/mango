(Disclaimer: Mango is extremely experimental, so don't whine to me when it eats your children.)

# Mango

Mango is a modular web-application framework inspired by [Rack](http://github.com/rack/rack) and [PEP333](http://www.python.org/dev/peps/pep-0333/).

## Overview

Mango is most of all a framework for other modules, and apps.  It takes a list of middleware, and an application and compiles them into a single http server object. The middleware and apps are written in a functional style, which keeps everything very self-contained.

Mango aims to make building reusable modules of HTTP functionality as easy as possible by enforcing a simple and unified API for web frameworks, web apps, middleware.

## API

As in Rack, the API is very minimal.

Applications are of the form:

  func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
    return 200, map[string]string{}, "Hello World!"
  }

Where:
  * mango.Env is a map[string]interface{} of the environment
    * It also has accessors for several other environment attributes:
      * mango.Env.Request() is the http.Request object
      * mango.Env.Session() is the map[string]interface for the session (only if using the Sessions middleware)
      * mango.Env.Logger() is the default logger for the app (or your custom logger if using the Logger middleware)
  * mango.Status is an integer for the HTTP status code for the response
  * mango.Headers is a map[string]string of the response headers
  * mango.Body is a string for the response body

## Available Modules

  * Sessions
    Usage: Sessions(app_secret, cookie_name string)
    Basic session management. Provides a mango.Env.Session() helper which returns a map[string]interface{} representing the session.  Any data stored in here will be serialized into the response session cookie.
    
  * Logger
    Usage: Logger(custom_logger *log.Logger)
    Provides a way to set a custom log.Logger object for the app. If this middleware is not provided Mango will set up a default logger to os.Stdout for the app to log to.

  * ShowErrors
    Usage: ShowErrors(templateString string)
    Catch any panics thrown from the app, and display them in an HTML template. If templateString is "", a default template is used. Not recommended to use in production as it could provide information helpful to attackers.

## Example App

  package main

  import (
    "mango"
  )

  func Hello(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
    env.Logger().Println("Got a", env.Request().Method, "request for", env.Request().RawUrl)
    return 200, map[string]string{}, mango.Body("Hello World!")
  }

  func main() {
    stack := new(mango.Stack)
    stack.Address = ":3000"
    stack.Middleware(ShowErrors(""))
    stack.Run(Hello)
  }


## Custom Middleware

Building middleware for Mango is moderately easy.

If you build some middleware and think others might find it useful, please let me know so I can include it in the core Mango source.  The success of Mango really depends on having excellent middleware packages.

An extremely basic middleware package is simply a function:

  func SilenceErrors(env mango.Env, app mango.App) (mango.Status, mango.Headers, mango.Body) {
    // Call our upstream app
    status, headers, body := app.Call(env)

    // If we got an error
    if status == 500 {
      // Silence it!
      status = 200
      headers = make(map[string]string)
      body = "Silence is golden!"
    }

    // Pass the response back to the client
    return status, headers, body
  }

To use this middleware we would do:

  func main() {
    stack := new(mango.Stack)
    stack.Address = ":3000"

    stack.Middleware(SilenceErrors) // Include our custom middleware

    stack.Run(Hello)
  }

For more complex middleware we may want to pass it configuration parameters. An example middleware package is one which will replace any image tags with funny pictures of cats:

  func Cats(cat\_images []string) mango.Middleware {
    // Initial setup stuff here
    // Done on application setup

    // Initialize our regex for finding image links
    regex := regexp.MustCompile("[\"']*(.jpg|.png)[\"']")

    // This is our middleware's request handler
    return func(env mango.Env, app mango.App) (mango.Status, mango.Headers, mango.Body) {
      // Call the upstream application
      status, headers, body := app(env)

      // Pick a random cat image
      image_url := cat_images[rand.Int()%len(cat_images)]

      // Substitute in our cat picture
      body = mango.Body(regex.ReplaceAllString(string(body), image_url))

      // Send the modified response onwards
      return status, headers, body
    }
  }

This works by building a closure (function) based on the parameters we pass, and returning it as the middleware. Through the magic of closures, the value we pass for cat\_images gets built into the function returned.

To use our middleware we would do:

  func main() {

    stack := new(mango.Stack)
    stack.Address = ":3000"

    // Initialize our cats middleware with our list of cat_images
    cat_images := []string{"ceiling_cat.jpg", "itteh_bitteh_kittehs.jpg", "monorail_cat.jpg"}
    cats_middleware = Cats(cat_images)

    stack.Middleware(cats_middleware) // Include the Cats middleware in our stack

    stack.Run(Hello)
  }


## About

Mango was written by [Paul Bellamy](http://paulbellamy.com). 

Follow me on [Twitter](http://www.twitter.com/pyrhho)!
