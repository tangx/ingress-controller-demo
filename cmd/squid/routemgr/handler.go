package routemgr

import (
	"fmt"
	"net/http"

	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

type MuxHandler struct {
	ReverseProxy *proxy.ReverseProxy
}

// 先不考虑多 backend 负载均衡的问题
func NewMuxHandler(server string, port int32) *MuxHandler {
	backend := fmt.Sprintf("%s:%d", server, port)
	return &MuxHandler{
		ReverseProxy: proxy.NewReverseProxy(backend),
	}
}

func (*MuxHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {}
