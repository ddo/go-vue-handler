package vue

import (
	"net/http"
	"path"
	"strings"
)

// Handler return a http.Handler that supports Vue Router app with history mode
func Handler(publicDir string) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		url := req.URL.String()

		// static files
		if strings.Contains(url, ".") || url == "/" {
			handler.ServeHTTP(w, req)
			return
		}

		// the all 404 gonna be served as root
		http.ServeFile(w, req, path.Join(publicDir, "/index.html"))
	})
}
