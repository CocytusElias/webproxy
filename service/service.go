package service

import (
	"github.com/gorilla/websocket"
	"webProxy/extern/logger"
	"webProxy/service/pubsub"
)

var conn *websocket.Conn

func Start() {
	logger.Info("pubSub starting...")
	pubsub.Start()
	logger.Info("transmit starting...")
	startTransmit()
	logger.Info("service starting...")
	startHttp()
}
