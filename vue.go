package vue

import (
	"net/http"
	"path"
)

// statusInterceptor wraps http.ResponseWriter to get the status code
// returned by http.FileServer handler. If the status is 404 (NotFound),
// do not send the response to the client.
type statusInterceptor struct {
	http.ResponseWriter
	status int
}

func (w *statusInterceptor) Write(p []byte) (int, error) {
	if w.status == http.StatusNotFound {
		return 0, nil
	}
	return w.ResponseWriter.Write(p)
}

func (w *statusInterceptor) WriteHeader(code int) {
	if code != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(code)
	}
	w.status = code
}

// Handler return a http.Handler that supports Vue Router app with history mode
func Handler(publicDir string) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		interceptor := &statusInterceptor{
			ResponseWriter: w,
		}
		handler.ServeHTTP(interceptor, req)
		if interceptor.status == http.StatusNotFound {
			w.WriteHeader(http.StatusOK)
			http.ServeFile(w, req, path.Join(publicDir, "index.html"))
		}
	})
}
