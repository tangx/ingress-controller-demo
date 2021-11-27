package routermgr

import (
	"fmt"
	"net/http"

	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
	netv1 "k8s.io/api/networking/v1"
)

// MuxHandler 满足 mux handler， 并包含 fasthttp resverse proxy
type MuxHandler struct {
	reverseProxy *proxy.ReverseProxy
	annotations  map[string]string
}

// 先不考虑多 backend 负载均衡的问题
func NewMuxHandler(backend netv1.IngressBackend, annotations map[string]string) *MuxHandler {
	// backend := fmt.Sprintf("%s:%d", server, port)
	server := fmt.Sprintf("%s:%d", backend.Service.Name, backend.Service.Port.Number)
	return &MuxHandler{
		reverseProxy: proxy.NewReverseProxy(server),
		annotations:  annotations,
	}
}

// ServeHTTP 满足 mux handler 的接口规则
func (*MuxHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {}

func (h *MuxHandler) ReverseProxy() *proxy.ReverseProxy {
	return h.reverseProxy
}
