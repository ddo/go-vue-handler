# go-vue-handler
Vue Router history mode with Golang

## installation

```bash
go get -u github.com/ddo/go-vue-handler
```

## usage

* build vue app to get ``index.html`` and ``dist`` folder
* serve it as a static folder with go server
* all the static files must has extension

## example

```
/
    public/
        dist
        index.html
    server.go
```


```go
package main

import (
	"net/http"

	"github.com/ddo/go-vue-handler"
)

const (
	port = "8080"
	publicDir = "./public"
)

func main() {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: vue.Handler(publicDir),
	}
	err := server.ListenAndServe()
	panic(err)
}
```