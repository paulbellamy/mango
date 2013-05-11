package mango

import (
	"bytes"
	"fmt"
	"net/http"
)

type MockResponseWriter struct {
  headers http.Header
  started bool
  Status  int
  Body    bytes.Buffer
}

func NewMockResponseWriter() *MockResponseWriter {
  return &MockResponseWriter{
    headers: make(map[string][]string),
  }
}

func (w *MockResponseWriter) Header() http.Header {
  return w.headers
}

func (w *MockResponseWriter) Write(p []byte) (int, error) {
  if !w.started {
    w.started = true
    w.Status = 200
  }
  return w.Body.Write(p)
}

func (w *MockResponseWriter) WriteHeader(status int) {
  if !w.started {
    w.started = true
    w.Status = status
  }
}

func init() {
	fmt.Println("Testing Mango Version:", VersionString())
}
