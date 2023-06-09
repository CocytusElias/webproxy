package main

import (
	"webProxy/client"
	"webProxy/client/config"
	"webProxy/client/module"
	"webProxy/extern/logger"
)

func init() {
	logger.Init("client")
	module.Init()
	config.Init()
}

// 首先链接 service 来创建 websocket
// 创建完成后，根据 websocket 拿到的数据来发起请求，并将请求结果返回给 service。

func main() {
	client.Start()
}
