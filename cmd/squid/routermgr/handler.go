package routermgr

import (
	"fmt"
	"net/http"

	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

// MuxHandler 满足 mux handler， 并包含 fasthttp resverse proxy
type MuxHandler struct {
	reverseProxy *proxy.ReverseProxy
}

// 先不考虑多 backend 负载均衡的问题
func NewMuxHandler(server string, port int32) *MuxHandler {
	backend := fmt.Sprintf("%s:%d", server, port)
	return &MuxHandler{
		reverseProxy: proxy.NewReverseProxy(backend),
	}
}

// ServeHTTP 满足 mux handler 的接口规则
func (*MuxHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {}

func (h *MuxHandler) ReverseProxy() *proxy.ReverseProxy {
	return h.reverseProxy
}
