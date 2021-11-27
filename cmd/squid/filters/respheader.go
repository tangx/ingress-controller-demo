package filters

import (
	"github.com/tangx/ingress-operator/cmd/squid/routermgr"
	"github.com/valyala/fasthttp"
)

type ResponseHeaderFilter struct {
}

func (rh *ResponseHeaderFilter) Do(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("age", "20")
}

func (rh *ResponseHeaderFilter) Type() routermgr.FilterType {
	return routermgr.FilterType_Response
}
