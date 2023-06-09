package utils

import (
	"go.uber.org/zap"
	"runtime/debug"
	"webProxy/extern/logger"
)

func SafeFunc(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			// 日志输出
			logger.Error("grountine panic", zap.Any("err", err))
			// 堆栈打印
			debug.PrintStack()
		}
	}()

	fn()
}
