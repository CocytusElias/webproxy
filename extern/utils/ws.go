package utils

import (
	"bytes"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"os"
	"webProxy/extern/constant"
	"webProxy/extern/logger"
)

// 图片缓存
var picMap = map[int][]byte{
	http.StatusRequestTimeout:      nil,
	http.StatusTooManyRequests:     nil,
	http.StatusInternalServerError: nil,
	http.StatusServiceUnavailable:  nil,
}

// 加载图片
func loadImage(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			logger.Error(err.Error())
		}
	}(file)

	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, file)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return buffer.Bytes(), nil
}

func GetWsResError(id int64, code int) *constant.WsRes {
	if code != http.StatusNotFound &&
		code != http.StatusRequestTimeout &&
		code != http.StatusTooManyRequests &&
		code != http.StatusInternalServerError &&
		code != http.StatusServiceUnavailable {
		code = http.StatusInternalServerError
	}

	var err error
	if picMap[code] == nil {
		if picMap[code], err = loadImage(fmt.Sprintf("static/%v.png", code)); err != nil {
			logger.Error(err.Error())
		}
	}

	return &constant.WsRes{
		ID:   id,
		Code: code,
		Header: map[string]string{
			"Content-Type": "image/png",
		},
		Body: picMap[code],
	}

}

func GetWsResByteError(id int64, code int) (resMessage []byte) {
	var err error
	if resMessage, err = sonic.Marshal(GetWsResError(id, code)); err == nil {
		logger.Error(err.Error())
	}
	return
}
