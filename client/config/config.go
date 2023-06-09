package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"net/url"
	"webProxy/client/module"
	"webProxy/extern/logger"
)

var WsServiceUrl string
var TimeoutSecond int64
var Routers []*Router

// 配置结构
type config struct {
	WsServiceUrl  string    `toml:"wsServiceUrl"`  // service ws 地址。
	TimeoutSecond int64     `toml:"timeoutSecond"` // 转发请求超时时间
	Routers       []*Router `toml:"routers"`       // 路由转发配置。
}

// Router 代理路由
type Router struct {
	Path           string           `toml:"path"`           // 路径匹配，正则。
	UpStream       string           `toml:"upStream"`       // 匹配路径转发服务地址。转发需要返回响应。
	CopyStreams    []string         `toml:"copyStreams"`    // 匹配路径拷贝转发服务地址。拷贝转发不需要响应。
	RewriteModules []*RewriteModule `toml:"rewriteModules"` // 路径重写模块配置。用数组格式来定义执行顺序
}

// RewriteModule 路由重写模块
type RewriteModule struct {
	Name   string `toml:"name"`   // 重写模块的模块名
	Params []any  `toml:"params"` // 重写模块传参
}

// Init 配置初始化
func Init() {

	// 读取并解析配置
	var conf *config
	if _, err := toml.DecodeFile("./config/client.toml", &conf); err != nil {
		panic(err)
	}

	// 配置挂载
	WsServiceUrl = conf.WsServiceUrl
	TimeoutSecond = conf.TimeoutSecond
	Routers = conf.Routers

	// ------ 启动前的配置检查 ------

	if TimeoutSecond <= 0 {
		logger.Panic("if the request forwarding timeout is set to a value less than or equal to 0, it may result in the inability to forward the request properly")
	}

	// 必须有转发路径才启动
	if Routers == nil || len(Routers) == 0 {
		logger.Panic("not find routers config")
	}

	// 所有转发路径也必须正常
	for idx, route := range Routers {
		// 既然有那就必须做转发配置
		if route == nil {
			logger.Panic(fmt.Sprintf("the path configuration is invalid. please fix it and try again. index location(start at 0): %v", idx))
		}

		// 请求后必须有响应，即使响应是空
		if upstreamAddr, err := url.Parse(route.UpStream); err != nil || upstreamAddr == nil || upstreamAddr.Host == "" {
			logger.Panic(fmt.Sprintf("in the current forwarding configuration, upstream is not valid, please check. index location(start at 0): %v", idx))
		}

		// 复制请求如果有，则做验证
		if route.CopyStreams != nil && len(route.CopyStreams) > 0 {
			for cdx, item := range route.CopyStreams {
				if copyUpstreamAddr, err := url.Parse(item); err != nil || copyUpstreamAddr == nil || copyUpstreamAddr.Host == "" {
					logger.Panic(fmt.Sprintf("in the current forwarding configuration, upstream is not valid, please check. index location(start at 0): %v, copy stream index location(start at 0): %v", idx, cdx))
				}
			}
		}

		// 重写模块如果有，则做验证
		if route.RewriteModules != nil && len(route.RewriteModules) > 0 {
			for mdx, item := range route.RewriteModules {
				if !module.Verify(item.Name, item.Params...) {
					logger.Panic(fmt.Sprintf("cannot find rewrite module. please fix it and try again. index location(start at 0): %v. rewrite module index location(start at 0): %v", idx, mdx))
				}
			}
		}

	}
}
