package main

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/tangx/ingress-operator/cmd/squid/config"
	"github.com/tangx/ingress-operator/cmd/squid/routermgr"
	"github.com/valyala/fasthttp"
)

func main() {

	cfg := config.NewConfig()
	cfg.Initial().ReadConfig()

	mgr := routermgr.NewRouterManager()
	mgr.ParseRules(cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("reverse proxy listen %s", addr)
	if err := fasthttp.ListenAndServe(
		addr,
		mgr.ProxyHandler,
	); err != nil {
		log.Fatal(err)
	}
}
