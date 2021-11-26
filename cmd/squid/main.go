package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tangx/ingress-operator/cmd/squid/config"
	"github.com/yeqown/log"

	"github.com/valyala/fasthttp"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

var (
	proxyServer  = proxy.NewReverseProxy("localhost:8080", proxy.WithTimeout(5*time.Second))
	proxyServer2 = proxy.NewReverseProxy("api-js.mixpanel.com")
	proxyServer3 = proxy.NewReverseProxy("www.jtthink.com")
)

func ProxyHandler(ctx *fasthttp.RequestCtx) {

	requestURI := string(ctx.RequestURI())
	log.Info("requestURI=", requestURI)

	if strings.HasPrefix(requestURI, "/local") {
		// "/local" path proxy to localhost
		arr := strings.Split(requestURI, "?")
		if len(arr) > 1 {
			arr = append([]string{"/foo"}, arr[1:]...)
			requestURI = strings.Join(arr, "?")
		}

		ctx.Request.SetRequestURI(requestURI)
		proxyServer.ServeHTTP(ctx)

		return
	}

	// 路由重写
	if strings.HasPrefix(requestURI, "/baidu") {

		// 删除 /baidu prefix。 类似于 rewrite
		// /baidu/abc/ -> /abc/
		newURI := strings.TrimLeft(requestURI, "/baidu")
		// 使用新的 uri
		ctx.Request.SetRequestURI(newURI)

		// 代理
		proxyServer3.ServeHTTP(ctx)

		return
	}

	// default 简单代理
	proxyServer2.ServeHTTP(ctx)
}

var cfg = config.NewConfig()

func main() {
	cfg.Initial().ReadConfig()

	logrus.Debugf("%+v", cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("reverse proxy listen %s", addr)
	if err := fasthttp.ListenAndServe(
		addr,
		ProxyHandler,
	); err != nil {
		log.Fatal(err)
	}
}
