package http

import (
	"io/fs"
	"net/http"
	"strings"
)

var FrontendFS fs.FS

func mountFrontend() {
	if FrontendFS == nil {
		http.HandleFunc("/", homeHandler)
		return
	}

	fileServer := http.FileServer(http.FS(FrontendFS))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/api/") {
			http.NotFound(w, r)
			return
		}

		_, err := fs.Stat(FrontendFS, strings.TrimPrefix(path, "/"))
		if err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		r.URL.Path = "/index.html"
		fileServer.ServeHTTP(w, r)
	})
}
