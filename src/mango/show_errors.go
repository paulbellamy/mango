package mango

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

func ShowErrors(templateString string, f http.HandlerFunc) http.HandlerFunc {
	if templateString == "" {
		templateString = `
      <html>
      <body>
        <p>
          {{.Error|html}}
        </p>
      </body>
      </html>
    `
	}

	errorTemplate := template.Must(template.New("error").Parse(templateString))

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buffer := &bytes.Buffer{}
				errorTemplate.Execute(buffer, struct{ Error string }{fmt.Sprintf("%s", err)})
				w.WriteHeader(500)
				w.Write(buffer.Bytes())
			}
		}()

		f(w, r)
	}
}
