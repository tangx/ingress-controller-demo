# 使用 fasthttp 实现反向代理

```go
package main

import (
	"github.com/yeqown/log"

	"github.com/valyala/fasthttp"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

var (
	proxyServer3 = proxy.NewReverseProxy("www.baidu.com")
)

func ProxyHandler(ctx *fasthttp.RequestCtx) {
	proxyServer3.ServeHTTP(ctx)
}

func main() {
	if err := fasthttp.ListenAndServe(":8081", ProxyHandler); err != nil {
		log.Fatal(err)
	}
}
```

```bash
curl http://localhost:8081/

<!DOCTYPE html>
<!--STATUS OK--><html> <head><meta http-equiv=content-type content=text/html;charset=utf-8><meta http-equiv=X-UA-Compatible content=IE=Edge><meta content=always name=referrer><link rel=stylesheet type=text/css href=http://s1.bdstatic.com/r/www/cache/bdorz/baidu.min.css><title>百度一下，你就知道</title></head> .... </html>

```