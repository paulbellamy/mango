package main

import (
	"fmt"
	"http"
	"log"
	"os"
	"time"
)

type MangoMiddleware struct {
	app  MangoApp
	Call MangoApp
}

type MangoApp func(map[string]interface{}) (int, map[string]string, string)

type Mango struct {
	address string
	app     MangoApp
}

func (this *Mango) BuildStack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env := make(map[string]interface{})
		status, headers, body := this.app(env)
		w.WriteHeader(status)
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		fmt.Fprintf(w, body)
	}
}

func (this *Mango) Run(app MangoApp) os.Error {
	this.app = app
	if this.address == "" {
		this.address = "0.0.0.0:8000"
	}

	log.Println("Starting Mango Server On:", this.address)
	http.HandleFunc("/", this.BuildStack())
	return http.ListenAndServe(this.address, nil)
}

func Hello(map[string]interface{}) (int, map[string]string, string) {
	return 200, map[string]string{"Content-Type": "text/html"}, fmt.Sprintf("%d", time.Seconds())
}

func main() {
	mango := new(Mango)
	mango.address = ":3000"
	mango.Run(Hello)
}
