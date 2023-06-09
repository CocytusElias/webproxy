package main

import (
	"webProxy/extern/logger"
	"webProxy/service"
	"webProxy/service/config"
)

func init() {
	logger.Init("service")
	config.Init()
}

// 首先创建 http 服务，用来接受请求
// 客户端向服务端发起请求并创建 socket 链接
// 后续 service 接受到的所有请求，都通过 websocket 转给 client。
// client 根据 service 接受到的请求后来发起请求，并将响应结果通过 websocket 返回给 service。
// service 内部建立了一套简单的发布订阅机制，service 接受到的请求发布进去，并订阅 client 返回的结果。

func main() {
	service.Start()
}
