package main

import (
	"fmt"
)

type MangoApp func(string) string
type MangoMiddleware func(string, []MangoMiddleware) string

func Logger(in string, upstream []MangoMiddleware) string {
	return fmt.Sprintf("Logger\n%s\nLogger", upstream[0](in, upstream[1:]))
}

func Sessions(in string, upstream []MangoMiddleware) string {
	return fmt.Sprintf("Sessions\n%s\nSessions", upstream[0](in, upstream[1:]))
}

func App(in string, upstream []MangoMiddleware) string {
	return "App"
}

func main() {
	fmt.Printf(Sessions("1", []MangoMiddleware{Logger, App}))
}
