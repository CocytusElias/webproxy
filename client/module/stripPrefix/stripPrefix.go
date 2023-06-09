package stripPrefix

import (
	"webProxy/extern/constant"
	"webProxy/extern/logger"
)

type Module struct {
}

// Init 初始化方法。
func Init() (module *Module, err error) {

	return nil, nil
}

// Handle 处理方法。
func (m *Module) Handle(wsReq *constant.WsReq, params ...any) (wsReqRewrite *constant.WsReqRewrite, err error) {
	strip := int(params[0].(int64))
	wsReqRewrite = &constant.WsReqRewrite{
		Method: wsReq.Method,
		Header: wsReq.Header,
		Body:   wsReq.Body,
		Path:   wsReq.Path,
	}

	if len(wsReqRewrite.Path) <= strip {
		wsReqRewrite.Path = ""
		logger.Warn("if you want to replace all characters, it is recommended to use stripAll instead of stripPrefix")
	} else {
		wsReqRewrite.Path = wsReqRewrite.Path[strip:]
	}

	return wsReqRewrite, nil
}

// Verify 参数验证方法。
func (m *Module) Verify(params ...any) bool {
	if params == nil || len(params) != 1 {
		logger.Error("stripPrefix requires and accepts only one argument")
		return false
	}

	if strip, ok := params[0].(int64); !ok || int(strip) <= 0 {
		logger.Error("the parameter for stripPrefix must be a positive integer")
		return false
	}

	return true
}
