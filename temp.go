package main

import (
	"fmt"
  "http"
	"log"
	"os"
	"time"
)

type MangoApp func(map[string]interface{}) (int, map[string]string, string)

type Mango struct {
	app MangoApp
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

func (this *Mango) App(app MangoApp) (os.Error) {
	this.app = app
  return nil
}

func (this *Mango) Run(address string) (os.Error) {
	log.Println("Starting Mango Server On:", address)
	http.HandleFunc("/", this.BuildStack())
	return http.ListenAndServe(address , nil)
}

func Hello(map[string]interface{}) (int, map[string]string, string) {
	return 200, map[string]string{"Content-Type": "text/html"}, fmt.Sprintf("%d", time.Seconds())
}

func main() {
	mango := new(Mango)
	mango.App(Hello)
	mango.Run("0.0.0.0:8080")
}
