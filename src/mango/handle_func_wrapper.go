package mango

import (
  "net/http"
)

type handlerFuncWrapper struct {
  f http.HandlerFunc
}

func (wrapper *handlerFuncWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  wrapper.f(w, r)
}

func FuncToHandler(f http.HandlerFunc) http.Handler {
  return &handlerFuncWrapper{f}
}

func HandlerToFunc(h http.Handler) http.HandlerFunc {
  return h.ServeHTTP
}
