package service

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"webProxy/extern/constant"
	"webProxy/extern/logger"
	"webProxy/service/config"
	"webProxy/service/pubsub"
)

// 启动 http 服务
func startHttp() {

	// 用于和 client 之间建立 websocket 链接通道
	http.HandleFunc("/proxy/chan", proxyChan)
	http.HandleFunc("/", proxyHttp)

	logger.Info(fmt.Sprintf("http service listen on: %s", config.Addr))

	err := http.ListenAndServe(config.Addr, nil)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

// 代理转发通道建立，使用 websocket 与代理客户端进行通信。因为涉及到客户端请求后返回。
func proxyChan(res http.ResponseWriter, req *http.Request) {
	var err error

	// ws 选项
	var upgrade = websocket.Upgrader{}

	// 将 HTTP 短链接升级为长链接
	if conn, err = upgrade.Upgrade(res, req, nil); err != nil {
		logger.Error(err.Error())
		return
	}
	defer func(conn *websocket.Conn) {
		if err = conn.Close(); err != nil {
			logger.Error(err.Error())
		}
	}(conn)

	// 监听读取消息
	for {
		var messageType int
		var message []byte
		if messageType, message, err = conn.ReadMessage(); err != nil {
			logger.Error(err.Error())
			if errors.Is(err, websocket.ErrCloseSent) {
				pubsub.PanicFullSub()
			}
			break
		}

		switch messageType {
		case websocket.TextMessage:

			// 解析文本数据并 sub
			var data constant.WsRes
			if err = sonic.Unmarshal(message, &data); err != nil {
				logger.Warn(err.Error(), zap.Int("messageType", messageType), zap.ByteString("message", message))
			}

			pubsub.Sub(&data)

		case websocket.PingMessage:
			if err = conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				return
			}
		case websocket.PongMessage:
			logger.Info("client pong success")
		case websocket.CloseMessage:
			pubsub.PanicFullSub()
			logger.Warn("websocket close")

		default:
			logger.Warn("unKnow message", zap.Int("messageType", messageType), zap.ByteString("message", message))
		}

	}

}

// http 请求建立
func proxyHttp(res http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(res, req.Body, 200<<20)
	//goland:noinspection ALL
	body, _ := ioutil.ReadAll(req.Body)

	pub := &constant.WsReq{
		Method: req.Method,
		Domain: req.Host,
		Path:   req.RequestURI,
		Header: map[string]string{},
		Body:   body,
	}

	if len(req.Header) > 0 {
		for k, v := range req.Header {
			pub.Header[k] = v[0]
		}
	}

	subChan := make(chan *constant.WsRes, 0)

	pubsub.Pub(pub, subChan)

	for data := range subChan {
		for k, v := range data.Header {
			res.Header().Set(k, v)
		}
		res.WriteHeader(data.Code)
		if _, err := res.Write(data.Body); err != nil {
			logger.Error(err.Error())
		}
		return
	}
}
