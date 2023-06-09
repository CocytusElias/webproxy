package service

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"webProxy/extern/constant"
	"webProxy/extern/logger"
	"webProxy/extern/utils"
	"webProxy/service/pubsub"
)

func startTransmit() {
	// 处理发布进通道的数据
	go utils.SafeFunc(func() {
		pubsub.PubChan(func(req *constant.WsReq) error {
			if conn == nil {
				return errors.New("the websocket connection has not been established")
			}

			// 格式化请求数据
			data, err := sonic.Marshal(req)
			if err != nil {
				logger.Error(err.Error(), zap.Int64("id", req.ID), zap.String("method", req.Method), zap.String("domain", req.Domain), zap.String("url", req.Path), zap.ByteString("body", req.Body))
				return err
			}

			// 尝试发送数据
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logger.Error(err.Error(), zap.Int64("id", req.ID), zap.String("method", req.Method), zap.String("domain", req.Domain), zap.String("url", req.Path), zap.ByteString("body", req.Body))
				return err
			}

			return nil
		})
	})

}
