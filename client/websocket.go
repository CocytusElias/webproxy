package client

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"time"
	"webProxy/client/config"
	"webProxy/client/transmit"
	"webProxy/extern/logger"
	"webProxy/extern/utils"
)

// 信号通道
var interrupt = make(chan os.Signal)
var conn *websocket.Conn

// 从 ws service 读取数据，并推入通道
func startWsClient() {
	var err error

	// 通过通道发送信号
	signal.Notify(interrupt, os.Interrupt)

	// 与 service 建立 websocket 链接
	conn, _, err = websocket.DefaultDialer.Dial(config.WsServiceUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	logger.Info("connecting to Websocket Server...")

	// 启动后台进程
	go utils.SafeFunc(background)

	// 会话保持
	go utils.SafeFunc(keepAlive)

	// 持续订阅并返回给 service
	go utils.SafeFunc(subMessage)

	// 监听读取消息
	for {
		var messageType int
		var message []byte
		if messageType, message, err = conn.ReadMessage(); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		switch messageType {
		case websocket.TextMessage:
			logger.Info("service message info", zap.Int("messageType", messageType), zap.ByteString("message", message))
			transmit.Pub(message)

		case websocket.PingMessage:
			if err = conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				return
			}
		case websocket.PongMessage:
			logger.Info("service pong success")
		case websocket.CloseMessage:
			logger.Warn("websocket close")

		default:
			logger.Warn("unKnow message", zap.Int("messageType", messageType), zap.ByteString("message", message))
		}
	}
}

// 会话保持
func keepAlive() {
	t := time.Tick(time.Second)

	for {
		<-t
		// 持续保持会话活跃
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			logger.Error(err.Error())
		}
	}
}

// 后台进程
func background() {
	for s := range interrupt {
		logger.Info(s.String())

		// Close our websocket connection
		err := conn.WriteMessage(websocket.CloseMessage, nil)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		select {
		case <-time.After(time.Second):
			logger.Info("Timeout in closing receiving channel. Exiting...")
		}
		return
	}
}

// 持续监听通道处理数据并返回给 service
func subMessage() {
	for message := range transmit.Sub() {
		for retry := 0; retry < 5; retry++ {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err == nil {
				break
			}
			logger.Error(err.Error())
		}
	}
}
