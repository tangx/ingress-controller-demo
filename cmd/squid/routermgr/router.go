package routermgr

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tangx/ingress-operator/cmd/squid/config"
	"github.com/tangx/ingress-operator/pkg/httpx"
	"github.com/valyala/fasthttp"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
	netv1 "k8s.io/api/networking/v1"
)

type RouterManager struct {
	*mux.Router
}

func NewRouterManager() *RouterManager {
	return &RouterManager{
		Router: mux.NewRouter(),
	}
}

func (mgr *RouterManager) ParseRules(cfg *config.Config) {
	for _, ing := range cfg.Ingresses {
		for _, rule := range ing.Rules {
			for _, path := range rule.HTTP.Paths {
				// 使用 path 创建 mux Route
				mgr.parsePath(path)
			}
		}
	}
}

func (mgr *RouterManager) parsePath(path netv1.HTTPIngressPath) {
	handler := NewMuxHandler(path.Backend.Service.Name, path.Backend.Service.Port.Number)

	// 创建 mux 路由， 并绑定 handler
	// 根据 path 类型创建不同的匹配方式
	switch mgr.pathType(path.PathType) {
	case netv1.PathTypeExact:
		mgr.NewRoute().Path(path.Path).Methods(httpx.MethodAny()...).Handler(handler)
	case netv1.PathTypeImplementationSpecific:
		// 使用下一条规则
		fallthrough
	default:
		// 默认为
		mgr.NewRoute().PathPrefix(path.Path).Methods(httpx.MethodAny()...).Handler(handler)
	}
}

// pathType 返回默认的 path type
func (mgr *RouterManager) pathType(typ *netv1.PathType) netv1.PathType {
	if typ == nil {
		return netv1.PathTypePrefix
	}

	return *typ
}

// GetReverseProxy 根据 fasthttp request 获取反代的 proxy handler
func (mgr *RouterManager) GetReverseProxy(req fasthttp.Request) *proxy.ReverseProxy {
	match := &mux.RouteMatch{}

	r := httpRequest(req)

	if mgr.Match(r, match) {
		return match.Handler.(*MuxHandler).ReverseProxy()
	}

	return nil
}

// httpRequest 根据 fasthttp request 创建 http request 用于进行路由匹配
func httpRequest(req fasthttp.Request) *http.Request {
	method := string(req.Header.Method())
	url := string(req.RequestURI())

	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil
	}
	return r
}

func (mgr *RouterManager) ProxyHandler(ctx *fasthttp.RequestCtx) {
	proxy := mgr.GetReverseProxy(ctx.Request)
	if proxy == nil {
		return
	}

	proxy.ServeHTTP(ctx)
}
