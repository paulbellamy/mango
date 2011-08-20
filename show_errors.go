package mango

import (
	"bytes"
	"template"
)

func ShowErrors(templateString string) Middleware {
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

	return func(env Env, app App) (status Status, headers Headers, body Body) {
		defer func() {
			if err := recover(); err != nil {
				buffer := bytes.NewBufferString("")
				errorTemplate.Execute(buffer, struct{ Error string }{err.(string)})
				status = 500
				headers = Headers{}
				body = Body(buffer.String())
			}
		}()

		return app(env)
	}
}
