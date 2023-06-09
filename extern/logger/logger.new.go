package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DateTimeEncoder is ...
const DateTimeEncoder = "2006-01-02 15:04:05"

// New 创建默认日志实例
func New(serverName string) *zap.Logger {
	return NewInfo(serverName)
}

//// NewDebug 创建标准输出为 Debug 等级的日志实例
//func NewDebug(serverName string) *zap.Logger {
//	return create(serverName, zapcore.DebugLevel)
//}

// NewInfo 创建标准输出为 Info 等级的日志实例
func NewInfo(serverName string) *zap.Logger {
	return create(serverName, zapcore.InfoLevel)
}

//// NewWarn 创建标准输出为 Warn 等级的日志实例
//func NewWarn(serverName string) *zap.Logger {
//	return create(serverName, zapcore.WarnLevel)
//}
//
//// NewError 创建标准输出为 Error 等级的日志实例
//func NewError(serverName string) *zap.Logger {
//	return create(serverName, zapcore.ErrorLevel)
//}

// create 创建日志
func create(serverName string, level zapcore.Level) *zap.Logger {
	if serverName == "" {
		serverName = "unKnow"
	}

	core := getLogCore(level)

	fields := zap.Fields(
		zap.String("serverName", serverName),
	)

	return zap.New(core, zap.AddCaller(), fields)
}

// 获取编码器配置
func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.CallerKey = "codeSite"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(DateTimeEncoder)
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	return encoderConfig
}

// 获取标准输出编码器
func getStdEncoder() zapcore.Encoder {
	encoderConfig := getEncoderConfig()

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // level大写染色编码器

	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		var encoder strings.Builder
		encoder.WriteString("[")
		encoder.WriteString(t.Format(DateTimeEncoder))
		encoder.WriteString("]")
		enc.AppendString(encoder.String())
	}

	encoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		var callerStr strings.Builder
		callerStr.WriteString("[")
		callerStr.WriteString(caller.TrimmedPath())
		callerStr.WriteString("]")
		enc.AppendString(callerStr.String())
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 获取日志内核
func getLogCore(level zapcore.Level) zapcore.Core {

	// 定义日志级别
	debugLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.DebugLevel
	})

	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.WarnLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	var zapCoreTee []zapcore.Core

	// 设置写入到标准输出
	stdWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))

	switch level {
	case zapcore.DebugLevel:
		zapCoreTee = append(zapCoreTee, zapcore.NewCore(getStdEncoder(), stdWriter, debugLevel))
		fallthrough
	case zapcore.InfoLevel:
		zapCoreTee = append(zapCoreTee, zapcore.NewCore(getStdEncoder(), stdWriter, infoLevel))
		fallthrough
	case zapcore.WarnLevel:
		zapCoreTee = append(zapCoreTee, zapcore.NewCore(getStdEncoder(), stdWriter, warnLevel))
		fallthrough
	case zapcore.ErrorLevel:
		fallthrough
	default:
		zapCoreTee = append(zapCoreTee, zapcore.NewCore(getStdEncoder(), stdWriter, errorLevel))
	}

	return zapcore.NewTee(zapCoreTee...)
}
