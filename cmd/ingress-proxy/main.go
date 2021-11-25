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
