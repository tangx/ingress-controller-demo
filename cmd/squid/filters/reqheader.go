package filters

import (
	"strings"

	"github.com/tangx/ingress-operator/cmd/squid/routermgr"
	"github.com/valyala/fasthttp"
)

const (
	ADD_REQUEST_HEADER = "squid.ingress.tangx.in/add-request-header"
)

type RequestHeaderFilter struct {
}

func (rh *RequestHeaderFilter) Do(ctx *fasthttp.RequestCtx, annotations map[string]string) {

	value, ok := annotations[ADD_REQUEST_HEADER]
	if !ok {
		return
	}

	for _, header := range strings.Split(value, ";") {
		if h := newHeader(header); h != nil {
			ctx.Request.Header.Add(h.key, h.value)
		}

	}
}

func (rh *RequestHeaderFilter) Type() routermgr.FilterType {
	return routermgr.FilterType_Request
}
