package main

import (
	"fmt"
)

type MangoMiddlewareInterface interface {
	Call(string) string
}

type MangoMiddleware struct {
	upstream MangoMiddlewareInterface
}

func (this *MangoMiddleware) Call(in string) string {
	return this.upstream.Call(in)
}


type Logger struct {
	MangoMiddleware
}

func (this *Logger) Call(in string) string {
	return fmt.Sprintf("Logger\n%s\nLogger", this.upstream.Call(in))
}

type Sessions struct {
	MangoMiddleware
}

func (this *Sessions) Call(in string) string {
	return fmt.Sprintf("Sessions\n%s\nSessions", this.upstream.Call(in))
}

type App struct{}

func (this *App) Call(in string) string {
	return "App"
}

func main() {
	logger := new(Logger)
	logger.upstream = new(App)
	sessions := new(Sessions)
	sessions.upstream = logger
	fmt.Printf(sessions.Call("1"))
}
