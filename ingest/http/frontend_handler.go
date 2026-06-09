package http

import (
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
)

var FrontendFS fs.FS

func mountFrontend() {
	if FrontendFS == nil {
		http.HandleFunc("/", homeHandler)
		return
	}

	indexHTML, err := fs.ReadFile(FrontendFS, "index.html")
	if err != nil {
		slog.Error("Failed to read embedded index.html", "error", err)
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

		fsPath := strings.TrimPrefix(path, "/")
		if fsPath == "" {
			fsPath = "index.html"
		}

		if fsPath != "index.html" {
			if _, err := fs.Stat(FrontendFS, fsPath); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
}
