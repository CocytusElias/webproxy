package transmit

import (
	"bytes"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
	"webProxy/client/config"
	"webProxy/client/module"
	"webProxy/extern/constant"
	"webProxy/extern/logger"
	"webProxy/extern/utils"
)

// 进行代理转发
func transmit(req *constant.WsReq) *constant.WsRes {
	router := getTransmitTarget(req)

	if router == nil {
		return utils.GetWsResError(req.ID, http.StatusNotFound)
	}

	// 请求重写处理
	if router.RewriteModules != nil && len(router.RewriteModules) > 0 {
		var err error
		if req, err = getHttpRequestInfo(req, router); err != nil {
			logger.Error(err.Error())
			return utils.GetWsResError(req.ID, http.StatusInternalServerError)
		}
	}

	// 发起拷贝请求，拷贝请求不管响应
	if router.CopyStreams != nil && len(router.CopyStreams) > 0 {
		for _, copyStream := range router.CopyStreams {
			go utils.SafeFunc(func() {
				statusCode, _, resHeader, err := requestStream(copyStream, req)
				if err != nil {
					logger.Error(err.Error(),
						zap.String("method", req.Method), zap.String("copyStream", copyStream), zap.String("path", req.Path),
						zap.Any("header", req.Header), zap.String("domain", req.Domain))
					return
				}

				logger.Info("copyStream request",
					zap.String("method", req.Method), zap.String("copyStream", copyStream), zap.String("path", req.Path),
					zap.Any("header", req.Header), zap.Int("statusCode", statusCode), zap.Any("resHeader", resHeader))
			})
		}
	}

	// 发起正常请求，正常请求的响应要给 service 的
	statusCode, resBody, resHeader, err := requestStream(router.UpStream, req)
	if err != nil {
		logger.Error(err.Error(),
			zap.String("method", req.Method), zap.String("host", router.UpStream), zap.String("path", req.Path),
			zap.Any("header", req.Header), zap.ByteString("body", req.Body), zap.String("domain", req.Domain))
		return utils.GetWsResError(req.ID, http.StatusInternalServerError)
	}

	return &constant.WsRes{
		ID:     req.ID,
		Code:   statusCode,
		Header: resHeader,
		Body:   resBody,
	}
}

// 获取代理目标
func getTransmitTarget(wsReq *constant.WsReq) *config.Router {
	for _, router := range config.Routers {
		if matched, _ := regexp.MatchString(router.Path, wsReq.Path); matched {
			return router
		}
	}
	return nil
}

// 获取请求信息重写
func getHttpRequestInfo(wsReq *constant.WsReq, router *config.Router) (reqRewrite *constant.WsReq, err error) {
	reqRewrite = wsReq

	for _, rewriteModule := range router.RewriteModules {
		var wsReqRewrite *constant.WsReqRewrite
		if wsReqRewrite, err = module.Handle(reqRewrite, rewriteModule.Name, rewriteModule.Params...); err != nil {
			logger.Error(err.Error())
			return
		}

		reqRewrite.Method = wsReqRewrite.Method
		reqRewrite.Path = wsReqRewrite.Path
		reqRewrite.Header = wsReqRewrite.Header
		reqRewrite.Body = wsReqRewrite.Body
	}

	return
}

// 发起请求
func requestStream(host string, wsReq *constant.WsReq) (statusCode int, resBody []byte, resHeader map[string]string, err error) {

	var addr *url.URL

	if addr, err = url.Parse(host); err != nil {
		logger.Error(err.Error())
		return
	}

	addr = addr.JoinPath(wsReq.Path)

	var requestUrl string
	if requestUrl, err = url.QueryUnescape(addr.String()); err != nil {
		logger.Error(err.Error())
		return
	}

	var req *http.Request
	if req, err = http.NewRequest(wsReq.Method, requestUrl, bytes.NewReader(wsReq.Body)); err != nil {
		logger.Error(err.Error())
		return
	}

	for k, v := range wsReq.Header {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: time.Duration(config.TimeoutSecond) * time.Second,
	}

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		logger.Error(err.Error())
		return
	}

	statusCode = res.StatusCode

	if len(res.Header) > 0 {
		resHeader = make(map[string]string, 0)
		for k, v := range res.Header {
			resHeader[k] = v[0]
		}
	}

	//goland:noinspection ALL
	resD, _ := ioutil.ReadAll(res.Body)
	resBody = resD

	return
}
