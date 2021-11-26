# 使用 mux 路由规则匹配创建 fasthttp ReverseProxy 反向代理规则

### 初始化
1. 创建 muxMatcher 对象， 管理所有 fasthttp ReverseProxy 
2. 解析 config，
    1. 获取 `/path` ， 创建 mux route 规则 `mux.NewRoute()`
    2. 通过 `backend` ， 创建 fasthttp ReverseProxy 规则。
    3. 想办法吧 mux route 和 fasthttp reverse proxy 关联起来。

### 请求处理
3. 在 fasthttp handler 中， 通过 ctx 获取 method 和 uri 请求
    1. 进行 muxMatcher 进行匹配
    2. 匹配成功， 使用 reverse proxy 进行代理




