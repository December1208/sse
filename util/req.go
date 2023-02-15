package util

import (
	"github.com/imroc/req"
	"go.uber.org/zap"
	"time"
)

func DoJsonRespPost(url string, respData interface{}, v ...interface{}) bool {
	req.SetTimeout(5 * time.Second)
	r, err := req.Post(url, v...)
	if err != nil {
		MyLogger.Error("fail to send verify msg", zap.Error(err))
		return false
	}
	resp := r.Response()
	if resp.StatusCode != 200 {
		resContent, _ := r.ToString()
		MyLogger.Warn("fail to send request",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("res content", resContent))
		return false
	}
	err = r.ToJSON(respData)
	if err != nil {
		MyLogger.Error("DoJsonRespPost decode json failed",
			zap.Error(err))
		return false
	}
	MyLogger.Info("DoJsonRespPost success", zap.Any("respData", respData))
	return true
}

func DoJsonRespGet(url string, respData interface{}, v ...interface{}) bool {
	req.SetTimeout(5 * time.Second)
	r, err := req.Get(url, v...)
	if err != nil {
		MyLogger.Error("fail to send verify msg", zap.Error(err))
		return false
	}
	resp := r.Response()
	if resp.StatusCode != 200 {
		resContent, _ := r.ToString()
		MyLogger.Warn("fail to send request",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("res content", resContent))
		return false
	}
	err = r.ToJSON(respData)
	if err != nil {
		MyLogger.Error("DoJsonRespPost decode json failed",
			zap.Error(err))
		return false
	}
	MyLogger.Info("DoJsonRespPost success", zap.Any("respData", respData))
	return true
}

func ReqGetWithRetry(retryTimes int, url string, respData interface{}, v ...interface{}) bool {
	for times := 0; times <= retryTimes; times++ {
		suc := DoJsonRespGet(url, respData, v)
		if suc {
			return true
		}
	}
	return false
}
