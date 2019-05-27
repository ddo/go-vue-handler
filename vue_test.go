package vue

import (
	"io/ioutil"
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

	tests := []struct {
		name     string
		path     string
		wantCode int
		wantBody []byte
	}{
		{"root", "/", http.StatusOK, indexhtml},
		{"index.html", "/index.html", http.StatusMovedPermanently, nil},
		{"js/app.js", "/js/app.js", http.StatusOK, appjs},
		{"foo", "/foo", http.StatusOK, indexhtml},
		{"bar", "/bar", http.StatusOK, indexhtml},
		{"foo.bar", "/foo.bar", http.StatusOK, indexhtml},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := Handler(publicDir)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantCode)
			}

			if got := rr.Body.Bytes(); !reflect.DeepEqual(got, tt.wantBody) {
				t.Error("handler returned wrong content when querying for " + tt.path)
			}
		})
	}
}
