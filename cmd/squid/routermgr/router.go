package routermgr

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tangx/ingress-operator/cmd/squid/config"
	"github.com/tangx/ingress-operator/pkg/httpx"
	"github.com/valyala/fasthttp"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
	netv1 "k8s.io/api/networking/v1"
)

type RouterManager struct {
	*mux.Router
	requestFilters  []ProxyFilter
	responseFilters []ProxyFilter
}

func NewRouterManager() *RouterManager {
	return &RouterManager{
		Router:          mux.NewRouter(),
		requestFilters:  make([]ProxyFilter, 0),
		responseFilters: make([]ProxyFilter, 0),
	}
}

func (mgr *RouterManager) ParseRules(cfg *config.Config) *RouterManager {
	for _, ing := range cfg.Ingresses {
		for _, rule := range ing.Rules {

			// 增加 mux host 验证
			route := mgr.NewRoute().Host(rule.Host)

			for _, path := range rule.HTTP.Paths {
				// 使用 path 创建 mux Route
				mgr.parsePath(route, path)
			}
		}
	}

	return mgr
}

func (mgr *RouterManager) parsePath(route *mux.Route, path netv1.HTTPIngressPath) {

	handler := NewMuxHandler(path.Backend.Service.Name, path.Backend.Service.Port.Number)

	// 创建 mux 路由， 并绑定 handler
	// 根据 path 类型创建不同的匹配方式
	switch mgr.pathType(path.PathType) {
	case netv1.PathTypeExact:
		route.Path(path.Path).Methods(httpx.MethodAny()...).Handler(handler)
	case netv1.PathTypeImplementationSpecific:
		// 使用下一条规则
		fallthrough
	default:
		// 默认为
		route.PathPrefix(path.Path).Methods(httpx.MethodAny()...).Handler(handler)
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
	url := fmt.Sprintf("%s://%s%s", req.Header.Method(), req.Host(), req.RequestURI())

	logrus.Debugf("url=>>> %s", url)
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

	for _, filter := range mgr.requestFilters {
		filter.Do(ctx)
	}

	proxy.ServeHTTP(ctx)

	for _, filter := range mgr.responseFilters {
		filter.Do(ctx)
	}
}

func (mgr *RouterManager) ProxyHandlerWithOptions(filters ...ProxyFilter) {
	for _, filter := range filters {
		switch filter.Type() {
		case FilterType_Request:
			mgr.requestFilters = append(mgr.requestFilters, filter)
		case FilterType_Response:
			mgr.responseFilters = append(mgr.responseFilters, filter)
		}
	}
}

type ProxyFilter interface {
	Type() FilterType
	Do(ctx *fasthttp.RequestCtx)
}

type FilterType string

const (
	FilterType_Request  FilterType = "RequestFilter"
	FilterType_Response FilterType = "ResponseFilter"
)
