package stripAll

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
	return &constant.WsReqRewrite{
		Method: wsReq.Method,
		Header: wsReq.Header,
		Body:   wsReq.Body,
		Path:   params[0].(string),
	}, nil
}

// Verify 参数验证方法。
func (m *Module) Verify(params ...any) bool {
	if params == nil || len(params) != 1 {
		logger.Error("stripAll requires and accepts only one argument")
		return false
	}

	if _, ok := params[0].(string); !ok {
		logger.Error("stripAll requires specifying a new path")
		return false
	}

	return true
}
