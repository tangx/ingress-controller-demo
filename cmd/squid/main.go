package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/tangx/ingress-operator/cmd/squid/config"
	"github.com/tangx/ingress-operator/cmd/squid/filters"
	"github.com/tangx/ingress-operator/cmd/squid/routermgr"
	"github.com/valyala/fasthttp"
)

func main() {

	cfg := config.NewConfig()
	cfg.Initial().ReadConfig()

	mgr := routermgr.NewRouterManager()
	mgr.ParseRules(cfg).
		WithFilters(
			&filters.RequestHeaderFilter{},
			&filters.ResponseHeaderFilter{},
		)

	// onWatch(mgr, cfg)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)

					cfg.ReadConfig()
					// mgr.Router = mux.NewRouter()
					mgr.ParseRules(cfg)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	logrus.Error("start to watch ")
	err = watcher.Add("config.yml")
	if err != nil {
		logrus.Fatal(err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("reverse proxy listen %s", addr)
	if err := fasthttp.ListenAndServe(
		addr,
		mgr.ProxyHandler,
	); err != nil {
		log.Fatal(err)
	}
}

func init() {
	if os.Getenv("env") == "local" {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
