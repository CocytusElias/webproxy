package client

import (
	"webProxy/client/transmit"
	"webProxy/extern/logger"
)

func Start() {
	logger.Info("transmit starting...")
	transmit.Start()
	logger.Info("client starting...")
	startWsClient()
}
