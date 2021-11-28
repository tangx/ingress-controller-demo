package routermgr

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tangx/ingress-operator/cmd/squid/config"
	"github.com/tangx/ingress-operator/pkg/httpx"
	"github.com/valyala/fasthttp"
	netv1 "k8s.io/api/networking/v1"
)

type RouterManager struct {
	*mux.Router
	requestFilters  []IProxyFilter
	responseFilters []IProxyFilter
}

func NewRouterManager() *RouterManager {
	return &RouterManager{
		Router:          mux.NewRouter(),
		requestFilters:  make([]IProxyFilter, 0),
		responseFilters: make([]IProxyFilter, 0),
	}
}

// func (mgr *RouterManager) NewRouter() *RouterManager {
// 	mgr.Router = mux.NewRouter()
// 	return mgr
// }

func (mgr *RouterManager) ParseRules(cfg *config.Config) *RouterManager {
	mgr.Router = mux.NewRouter()
	for _, ing := range cfg.Ingresses {
		for _, rule := range ing.Spec.Rules {

			// 增加 mux host 验证
			route := mgr.NewRoute().Host(rule.Host)

			for _, path := range rule.HTTP.Paths {
				// 使用 path 创建 mux Route
				handler := NewMuxHandler(path.Backend, ing.Annotations)

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
		}
	}

	return mgr
}

// pathType 返回默认的 path type
func (mgr *RouterManager) pathType(typ *netv1.PathType) netv1.PathType {
	if typ == nil {
		return netv1.PathTypePrefix
	}

	return *typ
}

// GetMuxHandler 根据 fasthttp request 获取反代的 proxy handler
func (mgr *RouterManager) GetMuxHandler(req fasthttp.Request) *MuxHandler {

	match := &mux.RouteMatch{}

	r := httpRequest(req)

	if mgr.Match(r, match) {
		return match.Handler.(*MuxHandler)
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

	handler := mgr.GetMuxHandler(ctx.Request)
	if handler == nil {
		return
	}

	proxy := handler.ReverseProxy()
	if proxy == nil {
		return
	}

	for _, filter := range mgr.requestFilters {
		filter.Do(ctx, handler.annotations)
	}

	proxy.ServeHTTP(ctx)

	for _, filter := range mgr.responseFilters {
		filter.Do(ctx, handler.annotations)
	}
}

func (mgr *RouterManager) WithFilters(filters ...IProxyFilter) {
	for _, filter := range filters {
		switch filter.Type() {
		case FilterType_Request:
			mgr.requestFilters = append(mgr.requestFilters, filter)
		case FilterType_Response:
			mgr.responseFilters = append(mgr.responseFilters, filter)
		}
	}
}
