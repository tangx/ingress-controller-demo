package routemgr

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tangx/ingress-operator/cmd/squid/config"
	"github.com/tangx/ingress-operator/pkg/httpx"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

type RuleManager struct {
	*mux.Router
}

func NewRuleManager() *RuleManager {
	return &RuleManager{
		Router: mux.NewRouter(),
	}
}

func (mgr *RuleManager) ParseRules(cfg *config.Config) {
	for _, rule := range cfg.Ingresses.Rules {
		for _, path := range rule.HTTP.Paths {
			// 使用 path 创建 mux Route
			handler := NewMuxHandler(path.Backend.Service.Name, path.Backend.Service.Port.Number)

			mgr.NewRoute().Path(path.Path).Methods(httpx.MethodAny()...).Handler(handler)
		}
	}
}

func (mgr *RuleManager) GetReverseProxy(req *http.Request) *proxy.ReverseProxy {
	match := &mux.RouteMatch{}
	ret := mgr.Match(req, match)

	if ret {
		return match.Handler.(*MuxHandler).ReverseProxy
	}

	return nil
}
