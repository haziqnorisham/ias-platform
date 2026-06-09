//go:build embed

package main

import (
	"embed"
	"io/fs"
	"log/slog"

	ingest_http "ias/automation/ingest/http"
)

//go:embed frontend-dist
var embeddedFrontend embed.FS

func init() {
	slog.Info("Frontend embedding enabled (build tag: embed)")
	sub, err := fs.Sub(embeddedFrontend, "frontend-dist")
	if err != nil {
		slog.Error("Failed to sub-embed frontend-dist", "error", err)
		return
	}
	ingest_http.FrontendFS = sub
}
