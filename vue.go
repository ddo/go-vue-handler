package vue

import (
	"bytes"
	"mime"
	"net/http"
	"path"
	"strings"
)

const html5mime = "text/html"

// statusInterceptor wraps http.ResponseWriter to get the status code
// returned by http.FileServer handler. If the status is 404 (NotFound),
// do not send the response to the client, however save it into an internal buffer.
type statusInterceptor struct {
	http.ResponseWriter
	status  int
	headers http.Header
	body    bytes.Buffer
}

func (w *statusInterceptor) Write(p []byte) (int, error) {
	if w.status == http.StatusNotFound {
		m := w.Header()
		for k, v := range m {
			w.headers[k] = v
			delete(m, k)
		}
		w.body.Write(p)
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

// Flush writes the data from internal buffers to original ResponseWriter.
func (w *statusInterceptor) Flush(code int) {
	m := w.Header()
	for k, v := range w.headers {
		m[k] = v
		delete(w.headers, k)
	}
	w.ResponseWriter.WriteHeader(code)
	w.ResponseWriter.Write(w.body.Bytes())
	w.body.Reset()
}

// matchAcceptHeader implements basic mime type matching for HTTP Accept header.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept
func matchAcceptHeader(want, got string) bool {
	wantMain := strings.Split(want, "/")[0]
	vals := strings.Split(got, ",")
	for _, v := range vals {
		m, _, err := mime.ParseMediaType(v)
		if err != nil {
			continue
		}
		if m == want || m == "*/*" || m == wantMain+"/*" {
			return true
		}
	}
	return false
}

// Handler return a http.Handler that supports Vue Router app with history mode
func Handler(publicDir string) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		interceptor := &statusInterceptor{
			ResponseWriter: w,
			headers:        make(http.Header),
		}
		handler.ServeHTTP(interceptor, req)
		if interceptor.status == http.StatusNotFound {
			accept := req.Header.Get("Accept")
			if matchAcceptHeader(html5mime, accept) {
				w.WriteHeader(http.StatusOK)
				http.ServeFile(w, req, path.Join(publicDir, "index.html"))
			} else {
				interceptor.Flush(http.StatusNotFound)
			}
		}
	})
}
