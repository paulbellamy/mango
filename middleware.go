package main

import (
	"fmt"
)

// This is the core app the user has written
type MangoApp func(string) string

// These are pieces of middleware,
// which 'wrap' around the core MangoApp
// (and each other)
type MangoMiddleware func(string, MangoApp) string

// Bundle a given list of MangoMiddleware pieces into a MangoApp
func Bundle(r ...MangoMiddleware) MangoApp {
	if len(r) <= 1 {
    // Terminate the innermost piece of MangoMiddleware
    // Basically stops it from recursing any further.
    return func(input string) string {
      return r[0](input, func(string) string {
        return ""
      })
    }
	}
	return wrap(r[0], Bundle(r[1:]...))
}

// Attach a piece of MangoMiddleware to the outside
// of a MangoApp. This wraps the inner MangoApp
// inside the outer MangoMiddleware.
func wrap(middleware MangoMiddleware, app MangoApp) MangoApp {
	return func(input string) string {
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
  return func(input string, upstream MangoApp) string {
    return app(input)
  }
}

func Sessions(in string, app MangoApp) string {
	return fmt.Sprintf("Sessions\n%s\nSessions", app(in))
}

func Logger(in string, app MangoApp) string {
	return fmt.Sprintf("Logger\n%s\nLogger", app(in))
}

func Hello(in string) string {
	return "Hello World!"
}

func main() {
	app := Bundle(Sessions, Logger, middlewareify(Hello))
	fmt.Printf(app(""))
}
