# go-vue-handler
Vue Router history mode with Golang

## installation

```bash
go get -u github.com/ddo/go-vue-handler
```

## usage

* build vue app to get ``index.html`` and ``dist`` folder
* serve it as a static folder with go server

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

## caveat

> Your server will no longer report 404 errors as all not-found paths now serve up your `index.html`
> file. To get around the issue, you should implement a catch-all route within your Vue app to show
> a 404 page:
>
> ```js
> const router = new VueRouter({
>   mode: 'history',
>   routes: [
>     { path: '*', component: NotFoundComponent }
>   ]
> })
> ```
>

https://router.vuejs.org/guide/essentials/history-mode.html#caveat

In addition, `go-vue-handler` tries to analyze HTTP Accept header, so if the browser expects to see an HTML5 page (mime type is `text/html`) it will receive `index.html` contents, otherwise - 404 status code.