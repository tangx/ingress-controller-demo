package filters

import (
	"github.com/tangx/ingress-operator/cmd/squid/routermgr"
	"github.com/valyala/fasthttp"
)

type RequestHeaderFilter struct {
}

func (rh *RequestHeaderFilter) Do(ctx *fasthttp.RequestCtx) {

	ctx.Request.Header.Add("name", "tangx.in")
}

func (rh *RequestHeaderFilter) Type() routermgr.FilterType {
	return routermgr.FilterType_Request
}
