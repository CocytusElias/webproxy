package constant

// WsRes websocket 响应结构
type WsRes struct {
	ID     int64             // 请求 ID，用来标识请求并返回响应的
	Code   int               // 响应 Code
	Header map[string]string // 响应头
	Body   []byte            // 响应体
}

// WsReq websocket 请求转发结构
type WsReq struct {
	ID     int64             // 请求 ID，用来标识请求并返回响应的
	Method string            // 请求方法，用来给 client 发起请求的
	Domain string            // 请求域名，用来给 client 发起请求的
	Path   string            // 请求地址，用来给 client 发起请求的
	Body   []byte            // 请求体，用来给 client 发起请求的
	Header map[string]string // 请求头，用来给 client 发起请求的
}

// WsReqRewrite 请求重写结构
type WsReqRewrite struct {
	Method string            // 请求方法，用来给 client 发起请求的
	Path   string            // 请求地址，用来给 client 发起请求的
	Body   []byte            // 请求体，用来给 client 发起请求的
	Header map[string]string // 请求头，用来给 client 发起请求的
}
