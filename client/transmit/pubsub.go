package transmit

import (
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"net/http"
	"webProxy/extern/constant"
	"webProxy/extern/logger"
	"webProxy/extern/utils"
)

// 发布请求转发
var pubPool = make(chan *constant.WsReq, 100)

// 订阅请求转发
var subPool = make(chan []byte, 100)

// Start 建立计算协程池，多协程计算处理。最多同时执行 1000 个
func Start() {
	for i := 0; i < 1000; i++ {
		go utils.SafeFunc(func() {

			for req := range pubPool {
				// 处理并获取要返回给 service 的消息
				res := transmit(req)

				// 将消息返回给 service
				var resMessage []byte
				var err error
				if resMessage, err = sonic.Marshal(&res); err != nil {
					logger.Error(err.Error())
					resMessage = utils.GetWsResByteError(req.ID, http.StatusInternalServerError)
				}
				subPool <- resMessage

			}
		})
	}
}

func Pub(message []byte) {
	// 解析 service 发送的消息
	var err error
	var req constant.WsReq
	if err = sonic.Unmarshal(message, &req); err != nil {
		logger.Error(err.Error())
		return
	}

	if req.ID <= 0 || req.Path == "" || req.Domain == "" || (req.Method != http.MethodGet &&
		req.Method != http.MethodPost &&
		req.Method != http.MethodPut &&
		req.Method != http.MethodDelete &&
		req.Method != http.MethodOptions) {
		logger.Warn("unKnow pub", zap.Int64("id", req.ID), zap.String("method", req.Method), zap.String("domain", req.Domain), zap.String("url", req.Path), zap.Any("header", req.Header), zap.ByteString("body", req.Body))
		return
	}

	pubPool <- &req
}

func Sub() <-chan []byte {
	return subPool
}
