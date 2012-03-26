package mango

import (
	"io/ioutil"
	"os"
	"path"
)

func fileIsRegular(fi os.FileInfo) bool {
	return fi.Mode()&(os.ModeDir|os.ModeSymlink|os.ModeNamedPipe|os.ModeSocket|os.ModeDevice) == 0
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	} else if !fileIsRegular(info) {
		return false
	}

	return true
}

func readFile(filename string) (string, error) {
	body, err := ioutil.ReadFile(filename)
	return string(body), err
}

func Static(directory string) Middleware {
	return func(env Env, app App) (Status, Headers, Body) {
		// See if we can serve a file
		file := path.Join(directory, env.Request().URL.Path)
		if fileExists(file) && (env.Request().Method == "GET" || env.Request().Method == "HEAD") {
			if body, err := readFile(file); err == nil {
				mime_type := []string{MimeType(path.Ext(file), "application/octet-stream")}
				return 200, Headers{"Content-Type": mime_type}, Body(body)
			} else {
				panic(err)
			}
		}

		// No file found, pass on to app
		return app(env)
	}
}
