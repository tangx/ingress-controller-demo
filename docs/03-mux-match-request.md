# 使用 mux 匹配路由规则

`github.com/gorilla/mux` 库支持 **正则匹配** 等多种匹配模式

在 [/demo/mux-match/main.go](/demo/mux-match/main.go) 路由匹配结果


```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tangx/ingress-operator/pkg/httpx"
)

func main() {
	r := mux.NewRouter()

	// 两条路由
	r.NewRoute().Path("/").Methods("GET")
	r.NewRoute().Path(`/user/{id:\d+}`).Methods(httpx.MethodAny()...)

	// 创建 matcher
	match := &mux.RouteMatch{}

	for _, mock := range []struct {
		Method string
		URL    string
		Wanted bool
	}{
		{
			Method: http.MethodGet,
			URL:    "https://www.baidu.com/",
			Wanted: true,
		},
		{
			Method: http.MethodPost,
			URL:    "https://www.tangx.in/user/123",
			Wanted: true,
		},
		{
			Method: http.MethodPut,
			URL:    "https://www.tangx.in/user/123",
			Wanted: true,
		},
		{
			Method: http.MethodPost,
			URL:    "/",
			Wanted: false,
		},
	} {
		// 创建一个 请求
		req, _ := http.NewRequest(mock.Method, mock.URL, nil)
		// 匹配路由
		ret := r.Match(req, match)
		// 匹配结果
		fmt.Printf("wanted: %t, real: %t\n", mock.Wanted, ret)
	}
}
// wanted: true, real: true
// wanted: true, real: true
// wanted: true, real: true
// wanted: false, real: false
```