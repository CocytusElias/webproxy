// 请不要直接更改此文件，此文件由 client-gen.go 生成

package module

import (
	"errors"
	"webProxy/client/module/stripAll"
	"webProxy/client/module/stripPrefix"
	"webProxy/client/module/stripSuffix"

	"webProxy/extern/constant"
	"webProxy/extern/logger"
)

// 用于模块数据处理
type m interface {
	Handle(wsReq *constant.WsReq, params ...any) (wsReqRewrite *constant.WsReqRewrite, err error)
	Verify(params ...any) bool
}

var modules = map[string]m{}

// Init 初始化
func Init() {

	var err error

	if modules["stripAll"], err = stripAll.Init(); err != nil {
		logger.Panic(err.Error())
	}

	if modules["stripPrefix"], err = stripPrefix.Init(); err != nil {
		logger.Panic(err.Error())
	}

	if modules["stripSuffix"], err = stripSuffix.Init(); err != nil {
		logger.Panic(err.Error())
	}

}

// Verify 验证
func Verify(moduleName string, params ...any) (exist bool) {
	var rewriteModule m
	if rewriteModule, exist = modules[moduleName]; !exist {
		return
	}

	return rewriteModule.Verify(params...)
}

// Handle 处理函数
func Handle(wsReq *constant.WsReq, name string, params ...any) (wsReqRewrite *constant.WsReqRewrite, err error) {

	module, exist := modules[name]
	if !exist {
		err = errors.New("not found rewrite module")
		return
	}

	return module.Handle(wsReq, params...)
}
