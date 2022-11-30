package sse_server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"sse_demo/util"
)

func GetFullMonitorData(interval, span int) []interface{} {
	tubeUrl := viper.GetString("inner_server.tube_host") + "/flame/v1/home/monitor"
	tubeRespData := gin.H{}
	tubeReqData := map[string]int{"interval": interval, "span": span}
	util.ReqGetWithRetry(2, tubeUrl, &tubeRespData, tubeReqData)

	htUrl := viper.GetString("inner_server.flame_host") + "/ht/v1/home/monitor"
	htRespData := gin.H{}
	htReqData := map[string]int{"interval": interval, "span": span}
	util.ReqGetWithRetry(2, htUrl, &htRespData, htReqData)

	var result []interface{}
	result = append(result, tubeRespData["result"].([]interface{})...)
	result = append(result, htRespData["result"].([]interface{})...)
	return result
}
