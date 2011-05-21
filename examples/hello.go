package main

import (
  "time"
  "mango"
)

func Hello(env Env) (Status, Headers, Body) {
	return 200, map[string]string{"Never-Gonna": "Give you up!"}, Body(fmt.Sprintf("%d", time.Seconds()))
}

func main() {
	mango := new(Mango)
	mango.address = ":3000"
	mango.Run(Hello)
}
