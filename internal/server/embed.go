package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/*
var embeddedWeb embed.FS

// WebHandler returns an http.Handler that serves the embedded React 19 Governance Studio.
func WebHandler() http.Handler {
	distFS, err := fs.Sub(embeddedWeb, "dist")
	if err != nil {
		// Fallback minimal handler if dist has not been populated
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>governmesh Governance Studio</title></head>
<body style="font-family: sans-serif; background: #0f172a; color: #fff; padding: 2rem;">
  <h2>governmesh AI Governance Studio</h2>
  <p>The embedded studio assets are initializing. Access the REST API at <a href="/api/health" style="color: #38bdf8">/api/health</a>.</p>
</body>
</html>`))
		})
	}

	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// SPA routing: if requesting a route without file extension, serve index.html
		if !strings.Contains(r.URL.Path, ".") {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
