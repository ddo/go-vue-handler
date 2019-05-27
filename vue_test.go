package vue

import (
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"
)

const publicDir = "test/public"

func readFile(t *testing.T, filename string) []byte {
	b, err := ioutil.ReadFile(filepath.Join(publicDir, filename))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestHandler(t *testing.T) {
	indexhtml := readFile(t, "index.html")
	appjs := readFile(t, "js/app.js")
	notfound := []byte("404 page not found\n")

	tests := []struct {
		name         string
		path         string
		acceptHeader string
		wantCType    string
		wantCode     int
		wantBody     []byte
	}{
		{"root", "/", html5mime, html5mime, http.StatusOK, indexhtml},
		{"index.html", "/index.html", html5mime, "", http.StatusMovedPermanently, nil},
		{"js/app.js", "/js/app.js", "", "application/javascript", http.StatusOK, appjs},
		{"foo", "/foo", html5mime, html5mime, http.StatusOK, indexhtml},
		{"bar", "/bar", html5mime, html5mime, http.StatusOK, indexhtml},
		{"foo.bar", "/foo.bar", html5mime, html5mime, http.StatusOK, indexhtml},
		{"missing.js", "/missing.js", "", "text/plain", http.StatusNotFound, notfound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.acceptHeader != "" {
				req.Header["Accept"] = []string{tt.acceptHeader}
			}

			rr := httptest.NewRecorder()
			handler := Handler(publicDir)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantCode)
			}

			gotCType, _, err := mime.ParseMediaType(rr.Header().Get("Content-Type"))
			if gotCType != "" && err != nil {
				t.Errorf("handler returned wrong content type: %v", err)
			}

			if gotCType != tt.wantCType {
				t.Errorf("handler returned wrong content type: got %v want %v",
					gotCType, tt.wantCType)
			}

			if got := rr.Body.Bytes(); !reflect.DeepEqual(got, tt.wantBody) {
				t.Error("handler returned wrong content when querying for " + tt.path)
			}
		})
	}
}
