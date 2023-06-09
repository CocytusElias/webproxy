package config

import (
	"github.com/BurntSushi/toml"
	"webProxy/extern/logger"
)

var TimeoutSecond int64
var Addr string
var MaxRequest int
var MaxRequestChannel int

// 配置结构
type config struct {
	TimeoutSecond     int64  `toml:"timeoutSecond"`     // 转发请求超时时间，此时间建议大于 client 的 timeoutSecond
	Addr              string `toml:"addr"`              // 监听地址
	MaxRequest        int    `toml:"maxRequest"`        // 最大请求数，超过此请求直接返回 429
	MaxRequestChannel int    `toml:"maxRequestChannel"` // 请求通道内最大数量，超过此请求直接返回 429
}

// Init 配置初始化
func Init() {

	// 读取并解析配置
	var conf *config
	if _, err := toml.DecodeFile("./config/service.toml", &conf); err != nil {
		panic(err)
	}

	if conf.Addr == "" {
		conf.Addr = "0.0.0.0:8080"
	}

	// 配置挂载
	TimeoutSecond = conf.TimeoutSecond
	Addr = conf.Addr
	MaxRequest = conf.MaxRequest
	MaxRequestChannel = conf.MaxRequestChannel

	// ------ 启动前的配置检查 ------

	if TimeoutSecond <= 0 {
		logger.Panic("if the request forwarding timeout is set to a value less than or equal to 0, it may result in the inability to forward the request properly")
	}

	if MaxRequest <= 1000 {
		logger.Panic("setting MaxRequest to a value less than 1000 may result in frequent 429 errors")
	}

	if MaxRequestChannel <= 100 {
		logger.Panic("setting MaxRequestChannel to a value less than 100 may result in frequent 429 errors")
	}

	if MaxRequest <= MaxRequestChannel {
		logger.Warn("if MaxRequest is smaller than MaxRequestChannel, it may trigger a 429 exception prematurely")
	}

}
