package logger

import (
	"go.uber.org/zap"
)

//goland:noinspection GoUnusedGlobalVariable
var Debug func(msg string, fields ...zap.Field)
var Info func(msg string, fields ...zap.Field)
var Warn func(msg string, fields ...zap.Field)
var Error func(msg string, fields ...zap.Field)
var Panic func(msg string, fields ...zap.Field)
var Zap *zap.Logger

func Init(name string) {
	Zap = New(name)

	Debug = Zap.Debug
	Info = Zap.Info
	Warn = Zap.Warn
	Error = Zap.Error
	Panic = Zap.Panic
}
