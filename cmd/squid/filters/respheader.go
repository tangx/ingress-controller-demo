package filters

import (
	"strings"

	"github.com/tangx/ingress-operator/cmd/squid/routermgr"
	"github.com/valyala/fasthttp"
)

const (
	ADD_RESPONSE_HEADER = "squid.ingress.tangx.in/add-response-header"
)

type ResponseHeaderFilter struct {
}

func (rh *ResponseHeaderFilter) Do(ctx *fasthttp.RequestCtx, annotations map[string]string) {

	value, ok := annotations[ADD_RESPONSE_HEADER]
	if !ok {
		return
	}

	for _, header := range strings.Split(value, ";") {
		if h := newHeader(header); h != nil {
			ctx.Response.Header.Add(h.key, h.value)
		}

	}
}

func (rh *ResponseHeaderFilter) Type() routermgr.FilterType {
	return routermgr.FilterType_Response
}
