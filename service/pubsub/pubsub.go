package pubsub

import (
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"time"
	"webProxy/extern/constant"
	"webProxy/extern/logger"
	"webProxy/extern/utils"
	"webProxy/service/config"
)

// 整体的处理逻辑如下：
// 	 1. service 接收到外部请求，并认为有效后，发布此请求到队列中。
// 	 2. 发布队列排队将此请求发送给 client。
// 	 3. client 处理完此请求并返回响应后，此请求就会进入订阅中。
// 	 4. service 拿到订阅并返回。

type subInfo struct {
	subChan chan *constant.WsRes
	pubUnix int64
}

var subMap = make(map[int64]*subInfo) // 订阅通道

// pub 发布通道
var pubChan = make(chan *constant.WsReq, 3000)

func Start() {
	// 后台处理，请求时间超过 30 秒的直接返回
	go utils.SafeFunc(func() {
		t := time.Tick(time.Second)
		for {
			<-t
			for id, s := range subMap {
				if s.pubUnix+config.TimeoutSecond <= time.Now().Unix() {
					Sub(utils.GetWsResError(id, http.StatusRequestTimeout))
				}
			}
		}
	})

}

// Pub 创建发布信息
func Pub(pub *constant.WsReq, subChan chan *constant.WsRes) {
	// 随机生成 id
	for {
		id := generateId()
		if _, exist := subMap[id]; !exist {
			pub.ID = id
			break
		}
	}

	// 设置订阅通道
	subMap[pub.ID] = &subInfo{
		subChan: subChan,
		pubUnix: time.Now().Unix(),
	}

	if len(subMap) > config.MaxRequest {
		Sub(utils.GetWsResError(pub.ID, http.StatusTooManyRequests))
	} else if len(pubChan) >= config.MaxRequestChannel {
		Sub(utils.GetWsResError(pub.ID, http.StatusTooManyRequests))
	}

	// 发布进通道
	pubChan <- pub
}

// PubChan 拉取发布通道内数据并调用处理方法处理
func PubChan(fn func(req *constant.WsReq) error) {
	for pub := range pubChan {
		if err := fn(pub); err != nil {
			logger.Error(err.Error(), zap.Int64("id", pub.ID), zap.String("method", pub.Method), zap.String("domain", pub.Domain), zap.String("url", pub.Path), zap.ByteString("body", pub.Body))
			Sub(utils.GetWsResError(pub.ID, http.StatusInternalServerError))
		}
	}
}

// Sub 发布「订阅数据」。
func Sub(sub *constant.WsRes) {

	if sub.ID <= 0 || sub.Code < 200 || sub.Code >= 600 || sub.Body == nil {
		logger.Warn("unKnow sub", zap.Int64("id", sub.ID), zap.Int("code", sub.Code), zap.Any("header", sub.Header))
		return
	}

	s, exist := subMap[sub.ID]
	delete(subMap, sub.ID)
	if !exist {
		logger.Warn("sub not found", zap.Int64("id", sub.ID), zap.Int("code", sub.Code), zap.Any("header", sub.Header))
		return
	}

	s.subChan <- sub

}

// PanicFullSub 与客户端的 socket 链接异常，所有的请求全部返回 500
func PanicFullSub() {
	subs := subMap
	subMap = make(map[int64]*subInfo)

	for id, s := range subs {
		s.subChan <- utils.GetWsResError(id, http.StatusServiceUnavailable)
	}

	subs = nil
}

// 生成请求 id
func generateId() int64 {
	// 获取当前的 Unix 时间戳（毫秒）
	now := time.Now().UnixNano() / int64(time.Millisecond)

	// 获取当前时间戳的后 8 位
	last8 := now % 100000000

	// 生成 6 位随机数
	rand.Seed(time.Now().UnixNano())
	random := rand.Int63n(1000000)

	// 合并后生成数字
	number := last8*1000000 + random

	return number
}
