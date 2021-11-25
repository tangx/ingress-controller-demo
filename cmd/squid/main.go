package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tangx/ingress-operator/cmd/squid/config"
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

var cfg = config.NewConfig()

func main() {
	cfg.Initial().ReadConfig()

	logrus.Debugf("%+v", cfg)

	if err := fasthttp.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Server.Port),
		ProxyHandler,
	); err != nil {
		log.Fatal(err)
	}
}
