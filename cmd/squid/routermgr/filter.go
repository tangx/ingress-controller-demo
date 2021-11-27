package routermgr

import "github.com/valyala/fasthttp"

type IProxyFilter interface {
	Type() FilterType
	Do(ctx *fasthttp.RequestCtx, annotation map[string]string)
}

type FilterType string

const (
	FilterType_Request  FilterType = "RequestFilter"
	FilterType_Response FilterType = "ResponseFilter"
)
