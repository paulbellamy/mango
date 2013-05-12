package mango

import (
	"net/http"
	"strings"
)

func fileExists(root http.FileSystem, path string) bool {
	if file, err := root.Open(path); err == nil {
		file.Close()
		return true
	}

	return false
}

func Static(root http.FileSystem, f http.HandlerFunc) http.HandlerFunc {
	fileServer := http.FileServer(root)

	return func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
		}

		// See if we can serve a file
		if fileExists(root, upath) && (r.Method == "GET" || r.Method == "HEAD") {
			fileServer.ServeHTTP(w, r)
		} else {
			// No file found, pass on to app
			if f != nil {
				f(w, r)
			}
		}
	}
}
